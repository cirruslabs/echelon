package console

import (
	"fmt"
	"github.com/cirruslabs/echelon/node"
	"os"
	"strings"
	"time"
)

type EchelonConsole struct {
	output       *os.File
	root         *node.EchelonNode
	currentFrame []string
	refreshRate  time.Duration
}

func NewConsole(output *os.File, root *node.EchelonNode) *EchelonConsole {
	return &EchelonConsole{
		output:      output,
		root:        root,
		refreshRate: 200 * time.Millisecond,
	}
}

func (console *EchelonConsole) StartDrawing() {
	for console.root != nil && console.root.IsRunning() {
		console.drawFrame()
		time.Sleep(console.refreshRate)
	}
	console.drawFrame()
}

func (console *EchelonConsole) drawFrame() {
	newFrame := console.root.Draw()
	sequenceToClearLineAndMoveUp := "\u001B[A\u001B[K" // http://ascii-table.com/ansi-escape-sequences.php
	_, _ = fmt.Fprint(console.output, strings.Repeat(sequenceToClearLineAndMoveUp, len(console.currentFrame)))
	for _, newLine := range newFrame {
		_, _ = fmt.Fprintln(console.output, newLine)
	}
	console.currentFrame = newFrame
}
