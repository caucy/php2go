package expr_test

import (
	"testing"

	"gotest.tools/assert"

	"github.com/i582/php2go/src/php/node/expr"
	"github.com/i582/php2go/src/php/position"

	"github.com/i582/php2go/src/php/node"
	"github.com/i582/php2go/src/php/node/stmt"
	"github.com/i582/php2go/src/php/php5"
	"github.com/i582/php2go/src/php/php7"
)

func TestInclude(t *testing.T) {
	src := `<? include $a;`

	expected := &node.Root{
		Position: &position.Position{
			StartLine: 1,
			EndLine:   1,
			StartPos:  3,
			EndPos:    14,
		},
		Stmts: []node.Node{
			&stmt.Expression{
				Position: &position.Position{
					StartLine: 1,
					EndLine:   1,
					StartPos:  3,
					EndPos:    14,
				},
				Expr: &expr.Include{
					Position: &position.Position{
						StartLine: 1,
						EndLine:   1,
						StartPos:  3,
						EndPos:    13,
					},
					Expr: &expr.Variable{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  11,
							EndPos:    13,
						},
						VarName: &node.Identifier{
							Position: &position.Position{
								StartLine: 1,
								EndLine:   1,
								StartPos:  11,
								EndPos:    13,
							},
							Value: "a",
						},
					},
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

func TestIncludeOnce(t *testing.T) {
	src := `<? include_once $a;`

	expected := &node.Root{
		Position: &position.Position{
			StartLine: 1,
			EndLine:   1,
			StartPos:  3,
			EndPos:    19,
		},
		Stmts: []node.Node{
			&stmt.Expression{
				Position: &position.Position{
					StartLine: 1,
					EndLine:   1,
					StartPos:  3,
					EndPos:    19,
				},
				Expr: &expr.IncludeOnce{
					Position: &position.Position{
						StartLine: 1,
						EndLine:   1,
						StartPos:  3,
						EndPos:    18,
					},
					Expr: &expr.Variable{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  16,
							EndPos:    18,
						},
						VarName: &node.Identifier{
							Position: &position.Position{
								StartLine: 1,
								EndLine:   1,
								StartPos:  16,
								EndPos:    18,
							},
							Value: "a",
						},
					},
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

func TestRequire(t *testing.T) {
	src := `<? require $a;`

	expected := &node.Root{
		Position: &position.Position{
			StartLine: 1,
			EndLine:   1,
			StartPos:  3,
			EndPos:    14,
		},
		Stmts: []node.Node{
			&stmt.Expression{
				Position: &position.Position{
					StartLine: 1,
					EndLine:   1,
					StartPos:  3,
					EndPos:    14,
				},
				Expr: &expr.Require{
					Position: &position.Position{
						StartLine: 1,
						EndLine:   1,
						StartPos:  3,
						EndPos:    13,
					},
					Expr: &expr.Variable{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  11,
							EndPos:    13,
						},
						VarName: &node.Identifier{
							Position: &position.Position{
								StartLine: 1,
								EndLine:   1,
								StartPos:  11,
								EndPos:    13,
							},
							Value: "a",
						},
					},
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

func TestRequireOnce(t *testing.T) {
	src := `<? require_once $a;`

	expected := &node.Root{
		Position: &position.Position{
			StartLine: 1,
			EndLine:   1,
			StartPos:  3,
			EndPos:    19,
		},
		Stmts: []node.Node{
			&stmt.Expression{
				Position: &position.Position{
					StartLine: 1,
					EndLine:   1,
					StartPos:  3,
					EndPos:    19,
				},
				Expr: &expr.RequireOnce{
					Position: &position.Position{
						StartLine: 1,
						EndLine:   1,
						StartPos:  3,
						EndPos:    18,
					},
					Expr: &expr.Variable{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  16,
							EndPos:    18,
						},
						VarName: &node.Identifier{
							Position: &position.Position{
								StartLine: 1,
								EndLine:   1,
								StartPos:  16,
								EndPos:    18,
							},
							Value: "a",
						},
					},
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
