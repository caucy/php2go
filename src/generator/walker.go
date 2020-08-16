package generator

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/i582/php2go/src/ctx"
	"github.com/i582/php2go/src/php/node"
	"github.com/i582/php2go/src/php/node/expr"
	"github.com/i582/php2go/src/php/node/expr/assign"
	"github.com/i582/php2go/src/php/node/expr/binary"
	"github.com/i582/php2go/src/php/node/name"
	"github.com/i582/php2go/src/php/node/scalar"
	"github.com/i582/php2go/src/php/node/stmt"
	"github.com/i582/php2go/src/php/walker"
	"github.com/i582/php2go/src/solver"
	"github.com/i582/php2go/src/types"
	"github.com/i582/php2go/src/utils"
	"github.com/i582/php2go/src/variable"
)

type GeneratorWalker struct {
	w        io.Writer
	filename string

	requireImports map[string]struct{}

	varStructDefinitionWriter *bytes.Buffer

	mainWriter   *bytes.Buffer
	headerWriter *bytes.Buffer

	varInfo types.VarInfo

	ctx *ctx.Context

	indents int
}

func NewGeneratorWalker(w io.Writer, filename string) GeneratorWalker {
	return GeneratorWalker{
		w:                         w,
		filename:                  filename,
		requireImports:            make(map[string]struct{}),
		mainWriter:                bytes.NewBufferString(""),
		headerWriter:              bytes.NewBufferString(""),
		varStructDefinitionWriter: bytes.NewBufferString(""),
		varInfo:                   types.NewVarInfo(),
	}
}

func (g GeneratorWalker) EnterChildNode(key string, w walker.Walkable) {}
func (g GeneratorWalker) LeaveChildNode(key string, w walker.Walkable) {}
func (g GeneratorWalker) EnterChildList(key string, w walker.Walkable) {}
func (g GeneratorWalker) LeaveChildList(key string, w walker.Walkable) {}
func (g *GeneratorWalker) LeaveNode(w walker.Walkable)                 {}

func (g *GeneratorWalker) EnterNode(w walker.Walkable) bool {
	n := w.(node.Node)

	switch n := n.(type) {
	case *node.Root:

	case *expr.ShortArray:
		return g.GenerateArray(n)
	case *expr.ArrayDimFetch:
		return g.GenerateArrayDimFetch(n)

	case *expr.FunctionCall:
		return g.GenerateFunctionCall(n)
	case *stmt.Function:
		return g.GenerateFunction(n)
	case *stmt.Return:
		return g.GenerateReturn(n)

	case *expr.Variable:
		return g.GenerateVariable(n)

	case *stmt.Expression:
		g.GenerateIndents()
		n.Expr.Walk(g)
		g.Write("\n")
		return false

	case *stmt.Echo:
		return g.GenerateEcho(n)
	case *assign.Assign:
		return g.GenerateAssign(n)

	case *stmt.For:
		return g.GenerateFor(n)
	case *stmt.While:
		return g.GenerateWhile(n)
	case *stmt.If:
		return g.GenerateIf(n)

	case *scalar.Lnumber:
		g.Write(n.Value)
	case *scalar.Dnumber:
		g.Write(n.Value)
	case *scalar.String:
		g.Write(n.Value)
	case *name.Name:
		val := utils.NamePartsToString(n.Parts)
		if val == "true" || val == "false" {
			g.Write(val)
		}

	case *binary.Plus:
		return g.GenerateBinaryOp(n)
	case *binary.Minus:
		return g.GenerateBinaryOp(n)
	case *binary.Mul:
		return g.GenerateBinaryOp(n)
	case *binary.Div:
		return g.GenerateBinaryOp(n)
	case *binary.Concat:
		return g.GenerateBinaryOp(n)

	case *binary.NotEqual:
		return g.GenerateBinaryOp(n)
	case *binary.Equal:
		return g.GenerateBinaryOp(n)
	case *binary.Smaller:
		return g.GenerateBinaryOp(n)
	case *binary.SmallerOrEqual:
		return g.GenerateBinaryOp(n)
	case *binary.Greater:
		return g.GenerateBinaryOp(n)
	case *binary.GreaterOrEqual:
		return g.GenerateBinaryOp(n)

	case *binary.BooleanAnd:
		return g.GenerateBinaryOp(n)
	case *binary.BooleanOr:
		return g.GenerateBinaryOp(n)

	case *expr.PostInc:
		n.Variable.Walk(g)
		g.Write("++")
		return false
	case *expr.PostDec:
		n.Variable.Walk(g)
		g.Write("--")
		return false
	case *expr.PreInc:
		g.Write("++")
		n.Variable.Walk(g)
		return false
	case *expr.PreDec:
		g.Write("--")
		n.Variable.Walk(g)
		return false
	}

	return true
}

func (g GeneratorWalker) Final() {
	_, _ = g.w.Write([]byte("// Code generated by php2go. PLEASE DO NOT EDIT.\n"))

	g.headerWriter.Write([]byte("package " + strings.TrimSuffix(g.filename, ".php") + "\n\n"))

	_, _ = g.w.Write(g.headerWriter.Bytes())

	if len(g.requireImports) != 0 {
		_, _ = g.w.Write([]byte("import (\n"))

		for imp := range g.requireImports {
			_, _ = g.w.Write([]byte(fmt.Sprintf("   \"%s\"\n", imp)))
		}

		_, _ = g.w.Write([]byte(")\n"))
	}

	_, _ = g.w.Write([]byte(g.varInfo.Generate()))

	_, _ = g.w.Write(g.mainWriter.Bytes())
}

func (g GeneratorWalker) WithNewContext() GeneratorWalker {
	newCtx := &ctx.Context{
		Variables:       variable.NewTable(),
		CurrentFunction: nil,
	}
	gg := g
	gg.ctx = newCtx
	return gg
}

func (g GeneratorWalker) WithContext(c *ctx.Context) GeneratorWalker {
	gg := g
	gg.ctx = c
	return gg
}

func (g GeneratorWalker) GenerateIndents() {
	for i := 0; i < g.indents; i++ {
		g.Write("\t")
	}
}

func (g *GeneratorWalker) Write(s string) {
	_, _ = g.mainWriter.Write([]byte(s))
}

func (g *GeneratorWalker) GenerateEcho(e *stmt.Echo) bool {
	g.GenerateIndents()
	g.requireImports["fmt"] = struct{}{}
	g.Write("fmt.Print(")
	g.ctx.InPrintFunctionCall = true
	for i, ex := range e.Exprs {
		ex.Walk(g)
		if i < len(e.Exprs)-1 {
			g.Write(", ")
		}
	}
	g.ctx.InPrintFunctionCall = false
	g.Write(")\n")
	return false
}

func (g *GeneratorWalker) GenerateArrayDimFetch(f *expr.ArrayDimFetch) bool {
	f.Variable.Walk(g)
	g.Write("[")
	f.Dim.Walk(g)
	g.Write("]")
	return false
}

func (g *GeneratorWalker) GenerateArray(a *expr.ShortArray) bool {
	if len(a.Items) == 0 {
		g.Write("[]Var{}")
		return false
	}

	isAssoc := a.Items[0].(*expr.ArrayItem).Key != nil
	for _, item := range a.Items {
		haveKey := item.(*expr.ArrayItem).Key != nil

		if isAssoc && !haveKey {
			panic("mixed array key")
		}
	}

	if isAssoc {
		return g.GenerateAssociativeArray(a)
	} else {
		return g.GeneratePlainArray(a)
	}

	return false
}

func (g *GeneratorWalker) GenerateAssociativeArray(a *expr.ShortArray) bool {
	valType := solver.ExprType(g.ctx, a.Items[0].(*expr.ArrayItem).Val)
	keyType := solver.ExprType(g.ctx, a.Items[0].(*expr.ArrayItem).Key)

	if !valType.SingleType() {
		g.Write("map[string]Var{")
	} else {
		g.Write("map[" + keyType.String() + "]" + valType.String() + "{")
		g.varInfo.AddTypes(types.NewTypes(types.NewAssociativeArrayType(keyType, valType, 1)))
	}

	for i, item := range a.Items {
		item := item.(*expr.ArrayItem)

		itemType := solver.ExprType(g.ctx, item.Val)
		if !valType.Equal(itemType) {
			panic("different types in array")
		}

		item.Key.Walk(g)

		g.Write(": ")

		item.Val.Walk(g)

		if i < len(a.Items)-1 {
			g.Write(", ")
		}
	}

	g.Write("}")

	return false
}

func (g *GeneratorWalker) GeneratePlainArray(a *expr.ShortArray) bool {
	valType := solver.ExprType(g.ctx, a.Items[0].(*expr.ArrayItem).Val)

	if !valType.SingleType() {
		g.Write("[]Var{")
	} else {
		g.Write("[]" + valType.String() + "{")
		g.varInfo.AddTypes(types.NewTypes(types.NewPlainArrayType(valType, 1)))
	}

	for i, item := range a.Items {
		item := item.(*expr.ArrayItem)

		itemType := solver.ExprType(g.ctx, item.Val)
		if !valType.Equal(itemType) {
			panic("different types in array")
		}

		item.Walk(g)

		if i < len(a.Items)-1 {
			g.Write(", ")
		}
	}

	g.Write("}")

	return false
}

func (g *GeneratorWalker) GenerateFor(f *stmt.For) bool {
	gg := g.WithContext(&f.Ctx)

	gg.GenerateIndents()
	gg.Write("for ")

	for i, init := range f.Init {
		init.Walk(&gg)
		if i < len(f.Init)-1 {
			gg.Write(", ")
		}
	}
	gg.Write("; ")

	for i, cond := range f.Cond {
		cond.Walk(&gg)
		if i < len(f.Cond)-1 {
			gg.Write(", ")
		}
	}
	gg.Write("; ")

	for i, aftereffect := range f.Loop {
		aftereffect.Walk(&gg)
		if i < len(f.Loop)-1 {
			gg.Write(", ")
		}
	}

	gg.Write(" {\n")
	gg.indents++

	f.Stmt.Walk(&gg)

	gg.indents--
	gg.GenerateIndents()
	gg.Write("}\n")

	return false
}

func (g *GeneratorWalker) GenerateWhile(wl *stmt.While) bool {
	gg := g.WithContext(&wl.Ctx)

	gg.GenerateIndents()
	gg.Write("for ")

	wl.Cond.Walk(&gg)

	gg.Write(" {\n")
	gg.indents++

	wl.Stmt.Walk(&gg)

	gg.indents--
	gg.GenerateIndents()
	gg.Write("}\n")

	return false
}

func (g *GeneratorWalker) GenerateIf(i *stmt.If) bool {
	gg := g.WithContext(&i.IfCtx)

	for _, v := range g.ctx.Variables.Vars {
		if v.FromIfElse && !v.WasInitialize {
			gg.GenerateIndents()
			gg.Write(fmt.Sprintf("var %s %s\n", v.Name, v.Type.GenerateName()))
			v.WasInitialize = true
		}
	}

	gg.GenerateIndents()
	gg.Write("if ")
	gg.ctx.InCondition = true
	i.Cond.Walk(&gg)
	gg.ctx.InCondition = false
	gg.Write(" {\n")
	gg.indents++

	i.Stmt.Walk(&gg)

	gg.indents--
	gg.GenerateIndents()
	gg.Write("}")

	if i.Else != nil {
		gg := g.WithContext(&i.ElseCtx)

		gg.Write(" else {\n")
		gg.indents++
		i.Else.Walk(&gg)
		gg.indents--
		gg.GenerateIndents()
		gg.Write("}\n")
	} else {
		gg.Write("\n")
	}

	return false
}

func (g *GeneratorWalker) GenerateBinaryOp(n node.Node) bool {
	switch n := n.(type) {
	case *binary.Plus:
		leftIsFloat := solver.ExprType(g.ctx, n.Left).Is(types.Float)
		rightIsFloat := solver.ExprType(g.ctx, n.Right).Is(types.Float)

		needLeftToFloat := !leftIsFloat && rightIsFloat
		needRightToFloat := leftIsFloat && !rightIsFloat

		if needLeftToFloat {
			g.Write("float64(")
		}

		n.Left.Walk(g)

		if needLeftToFloat {
			g.Write(")")
		}

		g.Write(" + ")

		if needRightToFloat {
			g.Write("float64(")
		}

		n.Right.Walk(g)

		if needRightToFloat {
			g.Write(")")
		}

	case *binary.Minus:
		leftIsFloat := solver.ExprType(g.ctx, n.Left).Is(types.Float)
		rightIsFloat := solver.ExprType(g.ctx, n.Right).Is(types.Float)

		needLeftToFloat := !leftIsFloat && rightIsFloat
		needRightToFloat := leftIsFloat && !rightIsFloat

		if needLeftToFloat {
			g.Write("float64(")
		}

		n.Left.Walk(g)

		if needLeftToFloat {
			g.Write(")")
		}

		g.Write(" - ")

		if needRightToFloat {
			g.Write("float64(")
		}

		n.Right.Walk(g)

		if needRightToFloat {
			g.Write(")")
		}

	case *binary.Mul:
		leftIsFloat := solver.ExprType(g.ctx, n.Left).Is(types.Float)
		rightIsFloat := solver.ExprType(g.ctx, n.Right).Is(types.Float)

		needLeftToFloat := !leftIsFloat && rightIsFloat
		needRightToFloat := leftIsFloat && !rightIsFloat

		if needLeftToFloat {
			g.Write("float64(")
		}

		n.Left.Walk(g)

		if needLeftToFloat {
			g.Write(")")
		}

		g.Write(" * ")

		if needRightToFloat {
			g.Write("float64(")
		}

		n.Right.Walk(g)

		if needRightToFloat {
			g.Write(")")
		}

	case *binary.Div:
		leftIsFloat := solver.ExprType(g.ctx, n.Left).Is(types.Float)
		rightIsFloat := solver.ExprType(g.ctx, n.Right).Is(types.Float)

		needLeftToFloat := !leftIsFloat && rightIsFloat
		needRightToFloat := leftIsFloat && !rightIsFloat

		if needLeftToFloat {
			g.Write("float64(")
		}

		n.Left.Walk(g)

		if needLeftToFloat {
			g.Write(")")
		}

		g.Write(" / ")

		if needRightToFloat {
			g.Write("float64(")
		}

		n.Right.Walk(g)

		if needRightToFloat {
			g.Write(")")
		}

	case *binary.Concat:
		n.Left.Walk(g)
		g.Write(" + ")
		n.Right.Walk(g)

	case *binary.Equal:
		g.ctx.InCompare = true
		leftType := solver.ExprType(g.ctx, n.Left)
		rightType := solver.ExprType(g.ctx, n.Right)
		if !leftType.SingleType() {
			n.Left.Walk(g)
			g.Write(fmt.Sprintf(".CompareWith%s(", utils.TransformType(rightType.String())))
			n.Right.Walk(g)
			g.Write(", Equal)")
		} else {
			n.Left.Walk(g)
			g.Write(" == ")
			n.Right.Walk(g)
		}
		g.ctx.InCompare = false

	case *binary.NotEqual:
		g.ctx.InCompare = true
		leftType := solver.ExprType(g.ctx, n.Left)
		rightType := solver.ExprType(g.ctx, n.Right)
		if !leftType.SingleType() {
			n.Left.Walk(g)
			g.Write(fmt.Sprintf(".CompareWith%s(", utils.TransformType(rightType.String())))
			n.Right.Walk(g)
			g.Write(", NotEqual)")
		} else {
			n.Left.Walk(g)
			g.Write(" != ")
			n.Right.Walk(g)
		}
		g.ctx.InCompare = false

	case *binary.Smaller:
		g.ctx.InCompare = true
		leftType := solver.ExprType(g.ctx, n.Left)
		rightType := solver.ExprType(g.ctx, n.Right)
		if !leftType.SingleType() {
			n.Left.Walk(g)
			g.Write(fmt.Sprintf(".CompareWith%s(", utils.TransformType(rightType.String())))
			n.Right.Walk(g)
			g.Write(", Smaller)")
		} else {
			n.Left.Walk(g)
			g.Write(" < ")
			n.Right.Walk(g)
		}
		g.ctx.InCompare = false

	case *binary.SmallerOrEqual:
		g.ctx.InCompare = true
		leftType := solver.ExprType(g.ctx, n.Left)
		rightType := solver.ExprType(g.ctx, n.Right)
		if !leftType.SingleType() {
			n.Left.Walk(g)
			g.Write(fmt.Sprintf(".CompareWith%s(", utils.TransformType(rightType.String())))
			n.Right.Walk(g)
			g.Write(", SmallerEqual)")
		} else {
			n.Left.Walk(g)
			g.Write(" < ")
			n.Right.Walk(g)
		}
		g.ctx.InCompare = false

	case *binary.Greater:
		g.ctx.InCompare = true
		leftType := solver.ExprType(g.ctx, n.Left)
		rightType := solver.ExprType(g.ctx, n.Right)
		if !leftType.SingleType() {
			n.Left.Walk(g)
			g.Write(fmt.Sprintf(".CompareWith%s(", utils.TransformType(rightType.String())))
			n.Right.Walk(g)
			g.Write(", Greater)")
		} else {
			n.Left.Walk(g)
			g.Write(" > ")
			n.Right.Walk(g)
		}
		g.ctx.InCompare = false

	case *binary.GreaterOrEqual:
		g.ctx.InCompare = true
		leftType := solver.ExprType(g.ctx, n.Left)
		rightType := solver.ExprType(g.ctx, n.Right)
		if !leftType.SingleType() {
			n.Left.Walk(g)
			g.Write(fmt.Sprintf(".CompareWith%s(", utils.TransformType(rightType.String())))
			n.Right.Walk(g)
			g.Write(", GreaterEqual)")
		} else {
			n.Left.Walk(g)
			g.Write(" >= ")
			n.Right.Walk(g)
		}
		g.ctx.InCompare = false

	case *binary.BooleanAnd:
		g.ctx.InBoolean = true
		n.Left.Walk(g)
		g.Write(" && ")
		g.ctx.InBoolean = true
		n.Right.Walk(g)
		g.ctx.InBoolean = false

	case *binary.BooleanOr:
		g.ctx.InBoolean = true
		n.Left.Walk(g)
		g.Write(" || ")
		g.ctx.InBoolean = true
		n.Right.Walk(g)
		g.ctx.InBoolean = false
	}

	return false
}

func (g *GeneratorWalker) GenerateReturn(r *stmt.Return) bool {
	g.GenerateIndents()
	g.Write("return ")

	tp := solver.ExprType(g.ctx, r.Expr)
	g.varInfo.AddTypes(tp)

	fn, need := g.ctx.CurrentFunction.ReturnType.GenerateCreation(tp)

	if need {
		g.Write(fn + "{ Val: ")
	}

	if r.Expr != nil {
		r.Expr.Walk(g)
	}

	if need {
		g.Write(" }")
	}

	g.Write("\n")

	return false
}

func (g *GeneratorWalker) GenerateFunctionCall(fn *expr.FunctionCall) bool {
	fnName := utils.NamePartsToString(fn.Function.(*name.Name).Parts)

	g.Write(fnName + "(")
	fn.ArgumentList.Walk(g)
	g.Write(")")

	return false
}

func (g *GeneratorWalker) GenerateAssign(a *assign.Assign) bool {
	e := a.Expression
	tp := solver.ExprType(g.ctx, e)

	switch a := a.Variable.(type) {
	case *expr.Variable:
		vr := a.Var
		if !vr.Type.Resolved() {
			vr.Type = solver.ResolveTypes(g.ctx, vr.Type)
		}

		if !vr.Type.ContainsMap(tp) {
			vr.Type.Merge(tp)
		}
		g.varInfo.AddTypes(vr.Type)

		vr.CurrentType = tp
		singleType := tp.SingleType()

		g.ctx.InAssign = true

		a.Walk(g)

		if vr.WasInitialize && !(singleType && !vr.Type.SingleType()) {
			g.Write(" = ")
		} else if !(singleType && !vr.Type.SingleType()) {
			g.Write(" := ")
			vr.WasInitialize = true
		}

		e.Walk(g)

		if singleType && !vr.Type.SingleType() {
			g.Write(")")
		}

		g.ctx.InAssign = false

	case *expr.ArrayDimFetch:
		isAddingElement := a.Dim == nil

		if isAddingElement {
			a.Variable.Walk(g)
			g.Write(" = append(")
			a.Variable.Walk(g)
			g.Write(", ")
			e.Walk(g)
			g.Write(")")
		} else {
			a.Variable.Walk(g)
			g.Write("[")
			a.Dim.Walk(g)
			g.Write("]")

			g.Write(" = ")

			e.Walk(g)
		}
	}

	return false
}

func (g *GeneratorWalker) GenerateVariable(v *expr.Variable) bool {
	if !v.Var.WasInitialize && !v.Var.Type.SingleType() {
		g.Write(v.Var.GenerateDefinition())
		v.Var.WasInitialize = true
		g.GenerateIndents()
	}

	g.varInfo.AddTypes(v.Var.Type)
	g.Write(v.Var.GenerateAccess(g.ctx.InAssign, g.ctx.InPrintFunctionCall, g.ctx.InCompare, g.ctx.InBoolean))

	return false
}

func (g *GeneratorWalker) GenerateFunction(f *stmt.Function) bool {
	g.ctx = &ctx.Context{
		Variables:       f.Func.Variables,
		CurrentFunction: f.Func,
	}

	if f.Func.ReturnType.Len() == 0 {
		g.Write(fmt.Sprintf("func %s() {\n", f.Func.Name))
	} else {
		if !f.Func.ReturnType.Resolved() {
			f.Func.ReturnType = solver.ResolveTypes(g.ctx, f.Func.ReturnType)
		}
		g.Write(fmt.Sprintf("func %s() %s {\n", f.Func.Name, f.Func.ReturnType.GenerateName()))
	}

	g.indents++

	for _, st := range f.Stmts {
		st.Walk(g)
	}

	g.Write("}\n\n")

	g.indents--

	return false
}