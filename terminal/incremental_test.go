package terminal_test

import (
	"bufio"
	"bytes"
	"github.com/cirruslabs/echelon/terminal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_calculateIncrementalUpdate_SameTwoLines(t *testing.T) {
	t.Parallel()
	var result bytes.Buffer
	terminal.CalculateIncrementalUpdate(
		bufio.NewWriter(&result),
		[]string{"Foo", "Bar"},
		[]string{"Foo", "Bar"},
	)
	assert.Equal(t, "", result.String())
}

func Test_calculateIncrementalUpdate_AddSingleLine(t *testing.T) {
	t.Parallel()
	var result bytes.Buffer
	terminal.CalculateIncrementalUpdate(
		bufio.NewWriter(&result),
		[]string{"Foo", "Bar"},
		[]string{"Foo", "Bar", "Baz"},
	)
	assert.Equal(t, "\rBaz\n", result.String())
}

func Test_calculateIncrementalUpdate_InplaceChange(t *testing.T) {
	t.Parallel()
	var result bytes.Buffer
	terminal.CalculateIncrementalUpdate(
		bufio.NewWriter(&result),
		[]string{"Foo", "Bar", "Baz"},
		[]string{"Foo", "Updated Bar", "Baz"},
	)
	assert.Equal(t, "\r\u001B[2A\u001B[KUpdated Bar\r\u001B[2B", result.String())
}
