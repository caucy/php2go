package expr

import (
	"github.com/i582/php2go/src/php/freefloating"
	"github.com/i582/php2go/src/php/node"
	"github.com/i582/php2go/src/php/position"
	"github.com/i582/php2go/src/php/walker"
)

// ConstFetch node
type ConstFetch struct {
	FreeFloating freefloating.Collection
	Position     *position.Position
	Constant     node.Node
}

// NewConstFetch node constructor
func NewConstFetch(Constant node.Node) *ConstFetch {
	return &ConstFetch{
		FreeFloating: nil,
		Constant:     Constant,
	}
}

// SetPosition sets node position
func (n *ConstFetch) SetPosition(p *position.Position) {
	n.Position = p
}

// GetPosition returns node positions
func (n *ConstFetch) GetPosition() *position.Position {
	return n.Position
}

func (n *ConstFetch) GetFreeFloating() *freefloating.Collection {
	return &n.FreeFloating
}

// Attributes returns node attributes as map
func (n *ConstFetch) Attributes() map[string]interface{} {
	return nil
}

// Walk traverses nodes
// Walk is invoked recursively until v.EnterNode returns true
func (n *ConstFetch) Walk(v walker.Visitor) {
	if v.EnterNode(n) == false {
		return
	}

	if n.Constant != nil {
		v.EnterChildNode("Constant", n)
		n.Constant.Walk(v)
		v.LeaveChildNode("Constant", n)
	}

	v.LeaveNode(n)
}
