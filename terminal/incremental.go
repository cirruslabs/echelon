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

func CalculateIncrementalUpdateMaxLines(output *bufio.Writer, linesBefore []string, linesAfter []string, maxLines int) {
	if len(linesBefore) > maxLines || len(linesAfter) > maxLines {
		linesToIgnoreBefore := len(linesBefore) - maxLines
		linesToIgnoreAfter := len(linesAfter) - maxLines
		linesToIgnore := linesToIgnoreBefore
		if linesToIgnore < linesToIgnoreAfter {
			linesToIgnore = linesToIgnoreAfter
		}
		linesBefore = linesBefore[linesToIgnore:]
		linesAfter = linesAfter[linesToIgnore:]
	}
	CalculateIncrementalUpdate(output, linesBefore, linesAfter)
}

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
	_, _ = output.WriteString(moveBeginningOfLine)

	if linesBeforeCount > 0 {
		// move up to the first line of the frame
		_, _ = output.WriteString(fmt.Sprintf("\x1B[%dA", linesBeforeCount))
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
					_, _ = output.WriteString(fmt.Sprintf("\x1B[%dB", linesSkipped))
				}
				_, _ = output.WriteString(eraseLine)
				_, _ = output.WriteString(linesAfter[i])
				_, _ = output.WriteString(moveBeginningOfLine)
				lastEditedIndex = i
			}
		}
		// in case last few lines were identical
		_, _ = output.WriteString(fmt.Sprintf("\x1B[%dB", linesMinCount-lastEditedIndex))
	}
	for i := linesMinCount; i < linesAfterCount; i++ {
		_, _ = output.WriteString(linesAfter[i])
		_, _ = output.WriteString("\n")
	}
	if linesBeforeCount > linesAfterCount {
		// erase everything down below
		_, _ = output.WriteString(eraseCursorDown)
	}
	_ = output.Flush()
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
