package terminal

import (
	"bufio"
	"fmt"
)

const (
	eraseLine           = "\x1B[K" // clear entire line
	eraseCursorDown     = "\x1B[J" // erase whole line
	moveBeginningOfLine = "\r"
)

func CalculateIncrementalUpdate(output *bufio.Writer, linesBefore []string, linesAfter []string) {
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
	output.WriteString(moveBeginningOfLine)

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
