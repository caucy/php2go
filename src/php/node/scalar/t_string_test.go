package scalar_test

import (
	"testing"

	"gotest.tools/assert"

	"github.com/i582/php2go/src/php/node"
	"github.com/i582/php2go/src/php/node/scalar"
	"github.com/i582/php2go/src/php/node/stmt"
	"github.com/i582/php2go/src/php/php5"
	"github.com/i582/php2go/src/php/php7"
	"github.com/i582/php2go/src/php/position"
)

func TestDoubleQuotedScalarString(t *testing.T) {
	src := `<? "test";`

	expected := &node.Root{
		Position: &position.Position{
			StartLine: 1,
			EndLine:   1,
			StartPos:  3,
			EndPos:    10,
		},
		Stmts: []node.Node{
			&stmt.Expression{
				Position: &position.Position{
					StartLine: 1,
					EndLine:   1,
					StartPos:  3,
					EndPos:    10,
				},
				Expr: &scalar.String{
					Position: &position.Position{
						StartLine: 1,
						EndLine:   1,
						StartPos:  3,
						EndPos:    9,
					},
					Value: "\"test\"",
				},
			},
		},
	}

	php7parser := php7.NewParser([]byte(src), "7.4")
	php7parser.Parse()
	actual := php7parser.GetRootNode()
	assert.DeepEqual(t, expected, actual)

	php5parser := php5.NewParser([]byte(src), "5.6")
	php5parser.Parse()
	actual = php5parser.GetRootNode()
	assert.DeepEqual(t, expected, actual)
}

func TestDoubleQuotedScalarStringWithEscapedVar(t *testing.T) {
	src := `<? "\$test";`

	expected := &node.Root{
		Position: &position.Position{
			StartLine: 1,
			EndLine:   1,
			StartPos:  3,
			EndPos:    12,
		},
		Stmts: []node.Node{
			&stmt.Expression{
				Position: &position.Position{
					StartLine: 1,
					EndLine:   1,
					StartPos:  3,
					EndPos:    12,
				},
				Expr: &scalar.String{
					Position: &position.Position{
						StartLine: 1,
						EndLine:   1,
						StartPos:  3,
						EndPos:    11,
					},
					Value: "\"\\$test\"",
				},
			},
		},
	}

	php7parser := php7.NewParser([]byte(src), "7.4")
	php7parser.Parse()
	actual := php7parser.GetRootNode()
	assert.DeepEqual(t, expected, actual)

	php5parser := php5.NewParser([]byte(src), "5.6")
	php5parser.Parse()
	actual = php5parser.GetRootNode()
	assert.DeepEqual(t, expected, actual)
}

func TestMultilineDoubleQuotedScalarString(t *testing.T) {
	src := `<? "
	test
	";`

	expected := &node.Root{
		Position: &position.Position{
			StartLine: 1,
			EndLine:   3,
			StartPos:  3,
			EndPos:    14,
		},
		Stmts: []node.Node{
			&stmt.Expression{
				Position: &position.Position{
					StartLine: 1,
					EndLine:   3,
					StartPos:  3,
					EndPos:    14,
				},
				Expr: &scalar.String{
					Position: &position.Position{
						StartLine: 1,
						EndLine:   3,
						StartPos:  3,
						EndPos:    13,
					},
					Value: "\"\n\ttest\n\t\"",
				},
			},
		},
	}

	php7parser := php7.NewParser([]byte(src), "7.4")
	php7parser.Parse()
	actual := php7parser.GetRootNode()
	assert.DeepEqual(t, expected, actual)

	php5parser := php5.NewParser([]byte(src), "5.6")
	php5parser.Parse()
	actual = php5parser.GetRootNode()
	assert.DeepEqual(t, expected, actual)
}

func TestSingleQuotedScalarString(t *testing.T) {
	src := `<? '$test';`

	expected := &node.Root{
		Position: &position.Position{
			StartLine: 1,
			EndLine:   1,
			StartPos:  3,
			EndPos:    11,
		},
		Stmts: []node.Node{
			&stmt.Expression{
				Position: &position.Position{
					StartLine: 1,
					EndLine:   1,
					StartPos:  3,
					EndPos:    11,
				},
				Expr: &scalar.String{
					Position: &position.Position{
						StartLine: 1,
						EndLine:   1,
						StartPos:  3,
						EndPos:    10,
					},
					Value: "'$test'",
				},
			},
		},
	}

	php7parser := php7.NewParser([]byte(src), "7.4")
	php7parser.Parse()
	actual := php7parser.GetRootNode()
	assert.DeepEqual(t, expected, actual)

	php5parser := php5.NewParser([]byte(src), "5.6")
	php5parser.Parse()
	actual = php5parser.GetRootNode()
	assert.DeepEqual(t, expected, actual)
}

func TestMultilineSingleQuotedScalarString(t *testing.T) {
	src := `<? '
	$test
	';`

	expected := &node.Root{
		Position: &position.Position{
			StartLine: 1,
			EndLine:   3,
			StartPos:  3,
			EndPos:    15,
		},
		Stmts: []node.Node{
			&stmt.Expression{
				Position: &position.Position{
					StartLine: 1,
					EndLine:   3,
					StartPos:  3,
					EndPos:    15,
				},
				Expr: &scalar.String{
					Position: &position.Position{
						StartLine: 1,
						EndLine:   3,
						StartPos:  3,
						EndPos:    14,
					},
					Value: "'\n\t$test\n\t'",
				},
			},
		},
	}

	php7parser := php7.NewParser([]byte(src), "7.4")
	php7parser.Parse()
	actual := php7parser.GetRootNode()
	assert.DeepEqual(t, expected, actual)

	php5parser := php5.NewParser([]byte(src), "5.6")
	php5parser.Parse()
	actual = php5parser.GetRootNode()
	assert.DeepEqual(t, expected, actual)
}
