package expr_test

import (
	"testing"

	"gotest.tools/assert"

	"github.com/i582/php2go/src/php/node/scalar"
	"github.com/i582/php2go/src/php/position"

	"github.com/i582/php2go/src/php/node/name"

	"github.com/i582/php2go/src/php/node/expr"

	"github.com/i582/php2go/src/php/node"
	"github.com/i582/php2go/src/php/node/expr/binary"
	"github.com/i582/php2go/src/php/node/stmt"
	"github.com/i582/php2go/src/php/php5"
	"github.com/i582/php2go/src/php/php7"
)

func TestFunctionCall(t *testing.T) {
	src := `<? foo();`

	expected := &node.Root{
		Position: &position.Position{
			StartLine: 1,
			EndLine:   1,
			StartPos:  3,
			EndPos:    9,
		},
		Stmts: []node.Node{
			&stmt.Expression{
				Position: &position.Position{
					StartLine: 1,
					EndLine:   1,
					StartPos:  3,
					EndPos:    9,
				},
				Expr: &expr.FunctionCall{
					Position: &position.Position{
						StartLine: 1,
						EndLine:   1,
						StartPos:  3,
						EndPos:    8,
					},
					Function: &name.Name{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  3,
							EndPos:    6,
						},
						Parts: []node.Node{
							&name.NamePart{
								Position: &position.Position{
									StartLine: 1,
									EndLine:   1,
									StartPos:  3,
									EndPos:    6,
								},
								Value: "foo",
							},
						},
					},
					ArgumentList: &node.ArgumentList{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  6,
							EndPos:    8,
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

func TestFunctionCallRelative(t *testing.T) {
	src := `<? namespace\foo();`

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
				Expr: &expr.FunctionCall{
					Position: &position.Position{
						StartLine: 1,
						EndLine:   1,
						StartPos:  3,
						EndPos:    18,
					},
					Function: &name.Relative{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  3,
							EndPos:    16,
						},
						Parts: []node.Node{
							&name.NamePart{
								Position: &position.Position{
									StartLine: 1,
									EndLine:   1,
									StartPos:  13,
									EndPos:    16,
								},
								Value: "foo",
							},
						},
					},
					ArgumentList: &node.ArgumentList{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  16,
							EndPos:    18,
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

func TestFunctionFullyQualified(t *testing.T) {
	src := `<? \foo([]);`

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
				Expr: &expr.FunctionCall{
					Position: &position.Position{
						StartLine: 1,
						EndLine:   1,
						StartPos:  3,
						EndPos:    11,
					},
					Function: &name.FullyQualified{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  3,
							EndPos:    7,
						},
						Parts: []node.Node{
							&name.NamePart{
								Position: &position.Position{
									StartLine: 1,
									EndLine:   1,
									StartPos:  4,
									EndPos:    7,
								},
								Value: "foo",
							},
						},
					},
					ArgumentList: &node.ArgumentList{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  7,
							EndPos:    11,
						},
						Arguments: []node.Node{
							&node.Argument{
								Position: &position.Position{
									StartLine: 1,
									EndLine:   1,
									StartPos:  8,
									EndPos:    10,
								},
								Variadic:    false,
								IsReference: false,
								Expr: &expr.ShortArray{
									Position: &position.Position{
										StartLine: 1,
										EndLine:   1,
										StartPos:  8,
										EndPos:    10,
									},
									Items: []node.Node{},
								},
							},
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

func TestFunctionCallVar(t *testing.T) {
	src := `<? $foo(yield $a);`

	expected := &node.Root{
		Position: &position.Position{
			StartLine: 1,
			EndLine:   1,
			StartPos:  3,
			EndPos:    18,
		},
		Stmts: []node.Node{
			&stmt.Expression{
				Position: &position.Position{
					StartLine: 1,
					EndLine:   1,
					StartPos:  3,
					EndPos:    18,
				},
				Expr: &expr.FunctionCall{
					Position: &position.Position{
						StartLine: 1,
						EndLine:   1,
						StartPos:  3,
						EndPos:    17,
					},
					Function: &expr.Variable{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  3,
							EndPos:    7,
						},
						VarName: &node.Identifier{
							Position: &position.Position{
								StartLine: 1,
								EndLine:   1,
								StartPos:  3,
								EndPos:    7,
							},
							Value: "foo",
						},
					},
					ArgumentList: &node.ArgumentList{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  7,
							EndPos:    17,
						},
						Arguments: []node.Node{
							&node.Argument{
								Position: &position.Position{
									StartLine: 1,
									EndLine:   1,
									StartPos:  8,
									EndPos:    16,
								},
								Variadic:    false,
								IsReference: false,
								Expr: &expr.Yield{
									Position: &position.Position{
										StartLine: 1,
										EndLine:   1,
										StartPos:  8,
										EndPos:    16,
									},
									Value: &expr.Variable{
										Position: &position.Position{
											StartLine: 1,
											EndLine:   1,
											StartPos:  14,
											EndPos:    16,
										},
										VarName: &node.Identifier{
											Position: &position.Position{
												StartLine: 1,
												EndLine:   1,
												StartPos:  14,
												EndPos:    16,
											},
											Value: "a",
										},
									},
								},
							},
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

func TestFunctionCallExprArg(t *testing.T) {
	src := `<? ceil($foo/3);`

	expected := &node.Root{
		Position: &position.Position{
			StartLine: 1,
			EndLine:   1,
			StartPos:  3,
			EndPos:    16,
		},
		Stmts: []node.Node{
			&stmt.Expression{
				Position: &position.Position{
					StartLine: 1,
					EndLine:   1,
					StartPos:  3,
					EndPos:    16,
				},
				Expr: &expr.FunctionCall{
					Position: &position.Position{
						StartLine: 1,
						EndLine:   1,
						StartPos:  3,
						EndPos:    15,
					},
					Function: &name.Name{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  3,
							EndPos:    7,
						},
						Parts: []node.Node{
							&name.NamePart{
								Position: &position.Position{
									StartLine: 1,
									EndLine:   1,
									StartPos:  3,
									EndPos:    7,
								},
								Value: "ceil",
							},
						},
					},
					ArgumentList: &node.ArgumentList{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  7,
							EndPos:    15,
						},
						Arguments: []node.Node{
							&node.Argument{
								Position: &position.Position{
									StartLine: 1,
									EndLine:   1,
									StartPos:  8,
									EndPos:    14,
								},
								Variadic:    false,
								IsReference: false,
								Expr: &binary.Div{
									Position: &position.Position{
										StartLine: 1,
										EndLine:   1,
										StartPos:  8,
										EndPos:    14,
									},
									Left: &expr.Variable{
										Position: &position.Position{
											StartLine: 1,
											EndLine:   1,
											StartPos:  8,
											EndPos:    12,
										},
										VarName: &node.Identifier{
											Position: &position.Position{
												StartLine: 1,
												EndLine:   1,
												StartPos:  8,
												EndPos:    12,
											},
											Value: "foo",
										},
									},
									Right: &scalar.Lnumber{
										Position: &position.Position{
											StartLine: 1,
											EndLine:   1,
											StartPos:  13,
											EndPos:    14,
										},
										Value: "3",
									},
								},
							},
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
