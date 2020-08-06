package node

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_AppendDescription(t *testing.T) {
	config := NewDefaultRenderingConfig()
	node := StartNewEchelonNode("test", config)
	node.Infof("1")
	node.Infof("2")
	node.Infof("3")
	node.Infof("4")
	node.Infof("5")
	node.Infof("6")
	assert.Equal(t, "123456", strings.Join(node.description, ""))
}
