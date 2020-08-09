package console

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_calculateIncrementalUpdate_SameTwoLines(t *testing.T) {
	var result bytes.Buffer
	calculateIncrementalUpdate(bufio.NewWriter(&result), []string{"Foo", "Bar"}, []string{"Foo", "Bar"})
	assert.Equal(t, "", result.String())
}

func Test_calculateIncrementalUpdate_AddSingleLine(t *testing.T) {
	var result bytes.Buffer
	calculateIncrementalUpdate(bufio.NewWriter(&result), []string{"Foo", "Bar"}, []string{"Foo", "Bar", "Baz"})
	assert.Equal(t, "Baz\n", result.String())
}

func Test_calculateIncrementalUpdate_InplaceChange(t *testing.T) {
	var result bytes.Buffer
	calculateIncrementalUpdate(bufio.NewWriter(&result), []string{"Foo", "Bar", "Baz"}, []string{"Foo", "Updated Bar", "Baz"})
	assert.Equal(t, "\u001B[A\u001B[A\u001B[KUpdated Bar\r\u001B[B\u001B[B", result.String())
}
