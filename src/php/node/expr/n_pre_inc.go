package expr

import (
	"github.com/i582/php2go/src/php/freefloating"
	"github.com/i582/php2go/src/php/node"
	"github.com/i582/php2go/src/php/position"
	"github.com/i582/php2go/src/php/walker"
)

// PreInc node
type PreInc struct {
	FreeFloating freefloating.Collection
	Position     *position.Position
	Variable     node.Node
}

// NewPreInc node constructor
func NewPreInc(Variable node.Node) *PreInc {
	return &PreInc{
		FreeFloating: nil,
		Variable:     Variable,
	}
}

// SetPosition sets node position
func (n *PreInc) SetPosition(p *position.Position) {
	n.Position = p
}

// GetPosition returns node positions
func (n *PreInc) GetPosition() *position.Position {
	return n.Position
}

func (n *PreInc) GetFreeFloating() *freefloating.Collection {
	return &n.FreeFloating
}

// Attributes returns node attributes as map
func (n *PreInc) Attributes() map[string]interface{} {
	return nil
}

// Walk traverses nodes
// Walk is invoked recursively until v.EnterNode returns true
func (n *PreInc) Walk(v walker.Visitor) {
	if v.EnterNode(n) == false {
		return
	}

	if n.Variable != nil {
		v.EnterChildNode("Variable", n)
		n.Variable.Walk(v)
		v.LeaveChildNode("Variable", n)
	}

	v.LeaveNode(n)
}
