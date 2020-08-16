package stmt_test

import (
	"testing"

	"gotest.tools/assert"

	"github.com/i582/php2go/src/php/node/expr/binary"
	"github.com/i582/php2go/src/php/position"

	"github.com/i582/php2go/src/php/node/expr"
	"github.com/i582/php2go/src/php/node/expr/assign"

	"github.com/i582/php2go/src/php/node/scalar"

	"github.com/i582/php2go/src/php/node"
	"github.com/i582/php2go/src/php/node/stmt"
	"github.com/i582/php2go/src/php/php5"
	"github.com/i582/php2go/src/php/php7"
)

func TestFor(t *testing.T) {
	src := `<? for($i = 0; $i < 10; $i++, $i++) {}`

	expected := &node.Root{
		Position: &position.Position{
			StartLine: 1,
			EndLine:   1,
			StartPos:  3,
			EndPos:    38,
		},
		Stmts: []node.Node{
			&stmt.For{
				Position: &position.Position{
					StartLine: 1,
					EndLine:   1,
					StartPos:  3,
					EndPos:    38,
				},
				Init: []node.Node{
					&assign.Assign{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  7,
							EndPos:    13,
						},
						Variable: &expr.Variable{
							Position: &position.Position{
								StartLine: 1,
								EndLine:   1,
								StartPos:  7,
								EndPos:    9,
							},
							VarName: &node.Identifier{
								Position: &position.Position{
									StartLine: 1,
									EndLine:   1,
									StartPos:  7,
									EndPos:    9,
								},
								Value: "i",
							},
						},
						Expression: &scalar.Lnumber{
							Position: &position.Position{
								StartLine: 1,
								EndLine:   1,
								StartPos:  12,
								EndPos:    13,
							},
							Value: "0",
						},
					},
				},
				Cond: []node.Node{
					&binary.Smaller{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  15,
							EndPos:    22,
						},
						Left: &expr.Variable{
							Position: &position.Position{
								StartLine: 1,
								EndLine:   1,
								StartPos:  15,
								EndPos:    17,
							},
							VarName: &node.Identifier{
								Position: &position.Position{
									StartLine: 1,
									EndLine:   1,
									StartPos:  15,
									EndPos:    17,
								},
								Value: "i",
							},
						},
						Right: &scalar.Lnumber{
							Position: &position.Position{
								StartLine: 1,
								EndLine:   1,
								StartPos:  20,
								EndPos:    22,
							},
							Value: "10",
						},
					},
				},
				Loop: []node.Node{
					&expr.PostInc{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  24,
							EndPos:    28,
						},
						Variable: &expr.Variable{
							Position: &position.Position{
								StartLine: 1,
								EndLine:   1,
								StartPos:  24,
								EndPos:    26,
							},
							VarName: &node.Identifier{
								Position: &position.Position{
									StartLine: 1,
									EndLine:   1,
									StartPos:  24,
									EndPos:    26,
								},
								Value: "i",
							},
						},
					},
					&expr.PostInc{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  30,
							EndPos:    34,
						},
						Variable: &expr.Variable{
							Position: &position.Position{
								StartLine: 1,
								EndLine:   1,
								StartPos:  30,
								EndPos:    32,
							},
							VarName: &node.Identifier{
								Position: &position.Position{
									StartLine: 1,
									EndLine:   1,
									StartPos:  30,
									EndPos:    32,
								},
								Value: "i",
							},
						},
					},
				},
				Stmt: &stmt.StmtList{
					Position: &position.Position{
						StartLine: 1,
						EndLine:   1,
						StartPos:  36,
						EndPos:    38,
					},
					Stmts: []node.Node{},
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

func TestAltFor(t *testing.T) {
	src := `<? for(; $i < 10; $i++) : endfor;`

	expected := &node.Root{
		Position: &position.Position{
			StartLine: 1,
			EndLine:   1,
			StartPos:  3,
			EndPos:    33,
		},
		Stmts: []node.Node{
			&stmt.AltFor{
				Position: &position.Position{
					StartLine: 1,
					EndLine:   1,
					StartPos:  3,
					EndPos:    33,
				},
				Cond: []node.Node{
					&binary.Smaller{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  9,
							EndPos:    16,
						},
						Left: &expr.Variable{
							Position: &position.Position{
								StartLine: 1,
								EndLine:   1,
								StartPos:  9,
								EndPos:    11,
							},
							VarName: &node.Identifier{
								Position: &position.Position{
									StartLine: 1,
									EndLine:   1,
									StartPos:  9,
									EndPos:    11,
								},
								Value: "i",
							},
						},
						Right: &scalar.Lnumber{
							Position: &position.Position{
								StartLine: 1,
								EndLine:   1,
								StartPos:  14,
								EndPos:    16,
							},
							Value: "10",
						},
					},
				},
				Loop: []node.Node{
					&expr.PostInc{
						Position: &position.Position{
							StartLine: 1,
							EndLine:   1,
							StartPos:  18,
							EndPos:    22,
						},
						Variable: &expr.Variable{
							Position: &position.Position{
								StartLine: 1,
								EndLine:   1,
								StartPos:  18,
								EndPos:    20,
							},
							VarName: &node.Identifier{
								Position: &position.Position{
									StartLine: 1,
									EndLine:   1,
									StartPos:  18,
									EndPos:    20,
								},
								Value: "i",
							},
						},
					},
				},
				Stmt: &stmt.StmtList{
					Position: &position.Position{
						StartLine: -1,
						EndLine:   -1,
						StartPos:  -1,
						EndPos:    -1,
					},
					Stmts: []node.Node{},
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
