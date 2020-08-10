package console

import (
	"bufio"
	"fmt"
	"github.com/cirruslabs/echelon/node"
	"io"
	"sync"
	"time"
)

const (
	defaultFrameBufSize = 38400    // 80 by 120 of 4 bytes UTF-8 characters
	eraseLine           = "\x1B[K" // clear entire line
	eraseCursorDown     = "\x1B[J" // erase whole line
	resetAutoWrap       = "\u001B[?7l"
	moveBeginningOfLine = "\r"
)

type EchelonConsole struct {
	output            *bufio.Writer
	nodes             []*node.EchelonNode
	currentFrameLines []string
	refreshRate       time.Duration
	renderRoot        bool
	drawLock          sync.Mutex
}

func NewConsole(output io.Writer, nodes []*node.EchelonNode) *EchelonConsole {
	return &EchelonConsole{
		output:      bufio.NewWriterSize(output, defaultFrameBufSize),
		nodes:       nodes,
		refreshRate: 200 * time.Millisecond,
	}
}

func (console *EchelonConsole) StartDrawing() {
	// don't wrap lines since it breaks incremental redraws
	console.output.WriteString(resetAutoWrap)
	for {
		if console.DrawFrame() {
			break
		}
		time.Sleep(console.refreshRate)
	}
}

func (console *EchelonConsole) DrawFrame() bool {
	console.drawLock.Lock()
	defer console.drawLock.Unlock()
	var newFrameLines []string
	var allCompleted = true
	for _, n := range console.nodes {
		if !n.HasCompleted() {
			allCompleted = false
		}
		newFrameLines = append(newFrameLines, n.Render()...)
	}
	oldFrame := console.currentFrameLines
	calculateIncrementalUpdate(console.output, oldFrame, newFrameLines)
	console.currentFrameLines = newFrameLines
	return allCompleted
}

func calculateIncrementalUpdate(output *bufio.Writer, linesBefore []string, linesAfter []string) {
	commonElements := commonElementsCount(linesBefore, linesAfter)
	if commonElements > 0 {
		linesBefore = linesBefore[commonElements:]
		linesAfter = linesAfter[commonElements:]
	}
	if len(linesBefore) == 0 && len(linesAfter) == 0 {
		// no changes
		return
	}
	linesBeforeCount := len(linesBefore)
	linesAfterCount := len(linesAfter)
	linesMinCount := linesBeforeCount
	if linesAfterCount < linesMinCount {
		linesMinCount = linesAfterCount
	}

	if linesBeforeCount > 0 {
		// move up to the first line of the frame
		output.WriteString(fmt.Sprintf("\x1B[%dA", linesBeforeCount))
	}
	if linesMinCount > 0 {
		// need to do incremental edits
		lastEditedIndex := 0
		for i := 0; i < linesMinCount; i++ {
			if linesBefore[i] != linesAfter[i] {
				// line existed before and was different so let's replace it
				linesSkipped := i - lastEditedIndex
				if linesSkipped > 0 {
					// move down
					output.WriteString(fmt.Sprintf("\x1B[%dB", linesSkipped))
				}
				output.WriteString(eraseLine)
				output.WriteString(linesAfter[i])
				output.WriteString(moveBeginningOfLine)
				lastEditedIndex = i
			}
		}
		// in case last few lines were identical
		output.WriteString(fmt.Sprintf("\x1B[%dB", linesMinCount-lastEditedIndex))
	}
	for i := linesMinCount; i < linesAfterCount; i++ {
		output.WriteString(linesAfter[i])
		output.WriteString("\n")
	}
	if linesBeforeCount > linesAfterCount {
		// erase everything down below
		output.WriteString(eraseCursorDown)
	}
	output.Flush()
}

func commonElementsCount(one []string, two []string) int {
	oneCount := len(one)
	twoCount := len(two)
	minCount := oneCount
	if twoCount < minCount {
		minCount = twoCount
	}
	for i := 0; i < minCount; i++ {
		if one[i] != two[i] {
			return i
		}
	}
	return minCount
}
