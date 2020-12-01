//nolint:testpackage
package renderers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_quotedIfNeeded(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "'foo'", quotedIfNeeded("foo"))

	assert.Equal(t, "task 'foo'", quotedIfNeeded("task 'foo'"))
	assert.Equal(t, "'foo' task", quotedIfNeeded("'foo' task"))
	assert.Equal(t, "task 'foo' has finished", quotedIfNeeded("task 'foo' has finished"))

	assert.Equal(t, "task \"foo\"", quotedIfNeeded("task \"foo\""))
	assert.Equal(t, "\"foo\" task", quotedIfNeeded("\"foo\" task"))
	assert.Equal(t, "task \"foo\" has finished", quotedIfNeeded("task \"foo\" has finished"))
}
