package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/i582/php2go/src/php/parser"
	"github.com/i582/php2go/src/php/printer"
	"github.com/i582/php2go/src/php/visitor"
	"github.com/pkg/profile"
	"github.com/yookoala/realpath"
)

var wg sync.WaitGroup
var phpVersion string
var dumpType string
var profiler string
var withFreeFloating *bool
var showResolvedNs *bool
var printBack *bool
var printPath *bool

type file struct {
	path    string
	content []byte
}

type result struct {
	path   string
	parser parser.Parser
}

func main() {
	withFreeFloating = flag.Bool("ff", false, "parse and show free floating strings")
	showResolvedNs = flag.Bool("r", false, "resolve names")
	printBack = flag.Bool("pb", false, "print AST back into the parsed file")
	printPath = flag.Bool("p", false, "print filepath")
	flag.StringVar(&dumpType, "d", "", "dump format: [custom, go, json, pretty_json]")
	flag.StringVar(&profiler, "prof", "", "start profiler: [cpu, mem, trace]")
	flag.StringVar(&phpVersion, "phpver", "7.4", "php version")

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		return
	}

	switch profiler {
	case "cpu":
		defer profile.Start(profile.ProfilePath("."), profile.NoShutdownHook).Stop()
	case "mem":
		defer profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()
	case "trace":
		defer profile.Start(profile.TraceProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()
	}

	numCpu := runtime.GOMAXPROCS(0)

	fileCh := make(chan *file, numCpu)
	resultCh := make(chan result, numCpu)

	// run 4 concurrent parserWorkers
	for i := 0; i < numCpu; i++ {
		go parserWorker(fileCh, resultCh)
	}

	// run printer goroutine
	go printerWorker(resultCh)

	// process files
	processPath(flag.Args(), fileCh)

	// wait the all files done
	wg.Wait()
	close(fileCh)
	close(resultCh)
}

func processPath(pathList []string, fileCh chan<- *file) {
	for _, path := range pathList {
		real, err := realpath.Realpath(path)
		checkErr(err)

		err = filepath.Walk(real, func(path string, f os.FileInfo, err error) error {
			if !f.IsDir() && filepath.Ext(path) == ".php" {
				wg.Add(1)
				content, err := ioutil.ReadFile(path)
				checkErr(err)
				fileCh <- &file{path, content}
			}
			return nil
		})
		checkErr(err)
	}
}

func parserWorker(fileCh <-chan *file, r chan<- result) {
	for {
		f, ok := <-fileCh
		if !ok {
			return
		}

		parserWorker, err := parser.NewParser(f.content, phpVersion)
		if err != nil {
			panic(err.Error())
		}

		if *withFreeFloating {
			parserWorker.WithFreeFloating()
		}

		parserWorker.Parse()

		r <- result{path: f.path, parser: parserWorker}
	}
}

func printerWorker(r <-chan result) {
	var counter int

	for {
		res, ok := <-r
		if !ok {
			return
		}

		counter++

		if *printPath {
			fmt.Fprintf(os.Stdout, "==> [%d] %s\n", counter, res.path)
		}

		for _, e := range res.parser.GetErrors() {
			fmt.Fprintf(os.Stdout, "==> %s\n", e)
		}

		if *printBack {
			o := bytes.NewBuffer([]byte{})
			p := printer.NewPrinter(o)
			p.Print(res.parser.GetRootNode())

			err := ioutil.WriteFile(res.path, o.Bytes(), 0644)
			checkErr(err)
		}

		var nsResolver *visitor.NamespaceResolver
		if *showResolvedNs {
			nsResolver = visitor.NewNamespaceResolver()
			res.parser.GetRootNode().Walk(nsResolver)
		}

		switch dumpType {
		case "custom":
			dumper := &visitor.Dumper{
				Writer:     os.Stdout,
				Indent:     "| ",
				NsResolver: nsResolver,
			}
			res.parser.GetRootNode().Walk(dumper)
		case "json":
			dumper := &visitor.JsonDumper{
				Writer:     os.Stdout,
				NsResolver: nsResolver,
			}
			res.parser.GetRootNode().Walk(dumper)
		case "pretty_json":
			dumper := &visitor.PrettyJsonDumper{
				Writer:     os.Stdout,
				NsResolver: nsResolver,
			}
			res.parser.GetRootNode().Walk(dumper)
		case "go":
			dumper := &visitor.GoDumper{Writer: os.Stdout}
			res.parser.GetRootNode().Walk(dumper)
		}

		wg.Done()
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
