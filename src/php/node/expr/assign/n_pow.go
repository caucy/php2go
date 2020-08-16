package assign

import (
	"github.com/i582/php2go/src/php/freefloating"
	"github.com/i582/php2go/src/php/node"
	"github.com/i582/php2go/src/php/position"
	"github.com/i582/php2go/src/php/walker"
)

// Pow node
type Pow struct {
	FreeFloating freefloating.Collection
	Position     *position.Position
	Variable     node.Node
	Expression   node.Node
}

// NewPow node constructor
func NewPow(Variable node.Node, Expression node.Node) *Pow {
	return &Pow{
		FreeFloating: nil,
		Variable:     Variable,
		Expression:   Expression,
	}
}

// SetPosition sets node position
func (n *Pow) SetPosition(p *position.Position) {
	n.Position = p
}

// GetPosition returns node positions
func (n *Pow) GetPosition() *position.Position {
	return n.Position
}

func (n *Pow) GetFreeFloating() *freefloating.Collection {
	return &n.FreeFloating
}

// Attributes returns node attributes as map
func (n *Pow) Attributes() map[string]interface{} {
	return nil
}

// Walk traverses nodes
// Walk is invoked recursively until v.EnterNode returns true
func (n *Pow) Walk(v walker.Visitor) {
	if v.EnterNode(n) == false {
		return
	}

	if n.Variable != nil {
		v.EnterChildNode("Variable", n)
		n.Variable.Walk(v)
		v.LeaveChildNode("Variable", n)
	}

	if n.Expression != nil {
		v.EnterChildNode("Expression", n)
		n.Expression.Walk(v)
		v.LeaveChildNode("Expression", n)
	}

	v.LeaveNode(n)
}
