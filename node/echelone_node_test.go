package node

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_AppendDescription(t *testing.T) {
	config := NewDefaultRenderingConfig()
	node := StartNewEchelonNode("test", config)
	node.AppendDescription("1")
	node.AppendDescription("2")
	node.AppendDescription("3")
	node.AppendDescription("4")
	node.AppendDescription("5")
	node.AppendDescription("6")
	assert.Equal(t, "123456", strings.Join(node.description, ""))
}
