package stmt

import (
	"github.com/i582/php2go/src/php/freefloating"
	"github.com/i582/php2go/src/php/node"
	"github.com/i582/php2go/src/php/position"
	"github.com/i582/php2go/src/php/walker"
)

// Unset node
type Unset struct {
	FreeFloating freefloating.Collection
	Position     *position.Position
	Vars         []node.Node
}

// NewUnset node constructor
func NewUnset(Vars []node.Node) *Unset {
	return &Unset{
		FreeFloating: nil,
		Vars:         Vars,
	}
}

// SetPosition sets node position
func (n *Unset) SetPosition(p *position.Position) {
	n.Position = p
}

// GetPosition returns node positions
func (n *Unset) GetPosition() *position.Position {
	return n.Position
}

func (n *Unset) GetFreeFloating() *freefloating.Collection {
	return &n.FreeFloating
}

// Attributes returns node attributes as map
func (n *Unset) Attributes() map[string]interface{} {
	return nil
}

// Walk traverses nodes
// Walk is invoked recursively until v.EnterNode returns true
func (n *Unset) Walk(v walker.Visitor) {
	if v.EnterNode(n) == false {
		return
	}

	if n.Vars != nil {
		v.EnterChildList("Vars", n)
		for _, nn := range n.Vars {
			if nn != nil {
				nn.Walk(v)
			}
		}
		v.LeaveChildList("Vars", n)
	}

	v.LeaveNode(n)
}
