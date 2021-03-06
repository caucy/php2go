package node_test

import (
	"testing"

	"gotest.tools/assert"

	"github.com/i582/php2go/src/php/position"
)

func TestPosition(t *testing.T) {
	expected := position.NewPosition(1, 1, 0, 1)
	for _, n := range nodes {
		n.SetPosition(expected)
		actual := n.GetPosition()
		assert.DeepEqual(t, expected, actual)
	}
}
