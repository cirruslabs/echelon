package node

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

// Reset ANSI sequence
const resetSequence = "\033[0m"

const (
	BLACK_COLOR = iota
	RED_COLOR
	GREEN_COLOR
	YELLOW_COLOR
	BLUE_COLOR
	MAGENTA_COLOR
	CYAN_COLOR
	WHITE_COLOR
)

type EchelonNode struct {
	lock                sync.RWMutex
	done                sync.WaitGroup
	title               string
	titleColor          int
	description         []string
	maxDescriptionLines int
	startTime           time.Time
	endTime             time.Time
	children            []*EchelonNode
}

func StartNewEchelonNode(title string) *EchelonNode {
	result := NewEchelonNode(title)
	result.Start()
	return result
}

func NewEchelonNode(title string) *EchelonNode {
	zeroTime := time.Time{}
	result := &EchelonNode{
		title:               title,
		titleColor:          -1,
		description:         make([]string, 0),
		maxDescriptionLines: 5,
		startTime:           zeroTime,
		endTime:             zeroTime,
		children:            make([]*EchelonNode, 0),
	}
	result.done.Add(1)
	return result
}

func (node *EchelonNode) UpdateTitle(text string) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.title = text
}

func (node *EchelonNode) ClearAllChildren() {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.children = make([]*EchelonNode, 0)
}

func (node *EchelonNode) ClearDescription() {
	node.SetDescription(make([]string, 0))
}

func (node *EchelonNode) SetDescription(description []string) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.description = description
}

func (node *EchelonNode) AppendDescription(line string) {
	node.AppendDescriptionLines([]string{line})
}

func (node *EchelonNode) AppendDescriptionLines(lines []string) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.description = append(node.description, lines...)
	linesTotal := len(node.description)
	if linesTotal > node.maxDescriptionLines {
		node.description = node.description[(linesTotal - node.maxDescriptionLines):]
	}
}

func (node *EchelonNode) AppendDescriptionBytes(bytes []byte) {
	node.lock.Lock()
	defer node.lock.Unlock()
	linesToAppend := strings.Split(string(bytes), "\n")
	if len(linesToAppend) == 0 {
		return
	}
	if len(node.description) == 0 {
		node.description = linesToAppend
		return
	}
	node.description[len(node.description)-1] = node.description[len(node.description)-1] + linesToAppend[0]
	if len(linesToAppend) > 1 {
		node.description = append(node.description, linesToAppend[1:]...)
	}
	linesTotal := len(node.description)
	if linesTotal > node.maxDescriptionLines {
		node.description = node.description[(linesTotal - node.maxDescriptionLines):]
	}
}

func (node *EchelonNode) Render() []string {
	node.lock.RLock()
	defer node.lock.RUnlock()
	result := []string{node.fancyTitle()}
	tail := node.description
	if len(node.children) > 0 {
		tail = node.RenderChildren()
	}
	for _, descriptionLine := range tail {
		result = append(result, "  "+descriptionLine)
	}

	return result
}

func (node *EchelonNode) RenderChildren() []string {
	node.lock.RLock()
	defer node.lock.RUnlock()
	var result []string
	for _, child := range node.children {
		result = append(result, child.Render()...)
	}
	return result
}

func (node *EchelonNode) fancyTitle() string {
	node.lock.RLock()
	defer node.lock.RUnlock()
	prefix := "[+]"
	if node.IsRunning() {
		prefix = "[-]"
	}
	coloredTitle := fmt.Sprintf("%s%s%s", getColorSequence(node.titleColor), node.title, resetSequence)
	return fmt.Sprintf("%s %s %s", prefix, coloredTitle, formatDuration(node.ExecutionDuration()))
}

func formatDuration(duration time.Duration) string {
	if duration < 10*time.Second {
		return fmt.Sprintf("%.1fs", float64(duration.Milliseconds())/1000)
	}
	if duration < time.Minute {
		return fmt.Sprintf("%ds", int(math.Floor(duration.Seconds()))%60)
	}
	if duration < time.Hour {
		return fmt.Sprintf("%02d:%02d", int(math.Floor(duration.Minutes()))%60, int(math.Floor(duration.Seconds()))%60)
	}
	return fmt.Sprintf(
		"%02d:%02d:%02d",
		int(math.Floor(duration.Hours())),
		int(math.Floor(duration.Minutes()))%60,
		int(math.Floor(duration.Seconds()))%60,
	)
}

func (node *EchelonNode) ExecutionDuration() time.Duration {
	node.lock.RLock()
	defer node.lock.RUnlock()
	if node.IsRunning() {
		return time.Now().Sub(node.startTime)
	} else {
		return node.endTime.Sub(node.startTime)
	}
}

func (node *EchelonNode) HasStarted() bool {
	node.lock.RLock()
	defer node.lock.RUnlock()
	return !node.startTime.IsZero()
}

func (node *EchelonNode) HasCompleted() bool {
	node.lock.RLock()
	defer node.lock.RUnlock()
	return !node.startTime.IsZero() && !node.endTime.IsZero()
}

func (node *EchelonNode) IsRunning() bool {
	node.lock.RLock()
	defer node.lock.RUnlock()
	return !node.startTime.IsZero() && node.endTime.IsZero()
}

func (node *EchelonNode) StartNewChild(childName string) *EchelonNode {
	child := StartNewEchelonNode(childName)
	node.AddNewChild(child)
	return child
}

func (node *EchelonNode) FindOrCreateChild(childTitle string) *EchelonNode {
	node.lock.Lock()
	defer node.lock.Unlock()
	// look from the end since this is a common pattern to get the last child
	for i := len(node.children) - 1; i >= 0; i-- {
		child := node.children[i]
		if child.title == childTitle {
			return child
		}
	}
	child := NewEchelonNode(childTitle)
	node.children = append(node.children, child)
	return child
}

func (node *EchelonNode) AddNewChild(child *EchelonNode) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.children = append(node.children, child)
}

func (node *EchelonNode) Start() {
	node.lock.Lock()
	defer node.lock.Unlock()
	if node.startTime.IsZero() {
		node.startTime = time.Now()
	}
}

func (node *EchelonNode) Complete() {
	node.CompleteWithColor(-1)
}
func (node *EchelonNode) CompleteWithColor(ansiColor int) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.titleColor = ansiColor
	if node.endTime.IsZero() {
		node.endTime = time.Now()
		node.done.Done()
	}
}

func (node *EchelonNode) WaitCompletion() {
	node.done.Wait()
}

func getColorSequence(code int) string {
	if code < 0 {
		return resetSequence
	}
	return fmt.Sprintf("\033[3%dm", code)
}
