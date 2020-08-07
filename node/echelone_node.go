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
	lock                    sync.RWMutex
	done                    sync.WaitGroup
	status                  string
	title                   string
	titleColor              int
	description             []string
	visibleDescriptionLines int
	config                  *EchelonNodeConfig
	startTime               time.Time
	endTime                 time.Time
	children                []*EchelonNode
}

func StartNewEchelonNode(title string, config *EchelonNodeConfig) *EchelonNode {
	result := NewEchelonNode(title, config)
	result.Start()
	return result
}

func NewEchelonNode(title string, config *EchelonNodeConfig) *EchelonNode {
	zeroTime := time.Time{}
	result := &EchelonNode{
		status:                  "⏸",
		title:                   title,
		titleColor:              -1,
		description:             make([]string, 0),
		visibleDescriptionLines: 5,
		config:                  config,
		startTime:               zeroTime,
		endTime:                 zeroTime,
		children:                make([]*EchelonNode, 0),
	}
	result.done.Add(1)
	return result
}

func (node *EchelonNode) UpdateTitle(text string) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.title = text
}

func (node *EchelonNode) UpdateConfig(config *EchelonNodeConfig) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.config = config
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

func (node *EchelonNode) SetVisibleDescriptionLines(count int) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.visibleDescriptionLines = count
}

func (node *EchelonNode) DescriptionLength() int {
	node.lock.RLock()
	defer node.lock.RUnlock()
	return len(node.description)
}

func (node *EchelonNode) Render() []string {
	node.lock.RLock()
	defer node.lock.RUnlock()
	result := []string{node.fancyTitle()}
	tail := node.renderChildren()
	if len(node.description) > node.visibleDescriptionLines {
		tail = append(tail, "...")
		tail = append(tail, node.description[(len(node.description)-node.visibleDescriptionLines):]...)
	} else {
		tail = append(tail, node.description...)
	}
	for _, descriptionLine := range tail {
		result = append(result, "  "+descriptionLine)
	}

	return result
}

func (node *EchelonNode) renderChildren() []string {
	var result []string
	for _, child := range node.children {
		result = append(result, child.Render()...)
	}
	return result
}

func (node *EchelonNode) fancyTitle() string {
	prefix := node.status
	if node.IsRunning() {
		prefix = node.config.CurrentProgressIndicatorFrame()
	}
	coloredTitle := node.title
	if node.titleColor >= 0 {
		coloredTitle = fmt.Sprintf("%s%s%s", getColorSequence(node.titleColor), node.title, resetSequence)
	}
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
	child := StartNewEchelonNode(childName, node.config)
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
	child := NewEchelonNode(childTitle, node.config)
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

func (node *EchelonNode) CompleteWithColor(status string, titleColor int) {
	if !node.endTime.IsZero() {
		return
	}
	node.lock.Lock()
	defer node.lock.Unlock()
	node.endTime = time.Now()
	node.status = status
	node.titleColor = titleColor
	node.done.Done()
}

func (node *EchelonNode) Complete() {
	if !node.endTime.IsZero() {
		return
	}
	node.lock.Lock()
	defer node.lock.Unlock()
	node.endTime = time.Now()
	node.done.Done()
}

func (node *EchelonNode) SetTitleColor(ansiColor int) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.titleColor = ansiColor
}

func (node *EchelonNode) SetStatus(text string) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.status = text
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

func (node *EchelonNode) Tracef(format string, args ...interface{}) {
	node.Logf(TraceLevel, format, args...)
}

func (node *EchelonNode) Debugf(format string, args ...interface{}) {
	node.Logf(DebugLevel, format, args...)
}

func (node *EchelonNode) Infof(format string, args ...interface{}) {
	node.Logf(InfoLevel, format, args...)
}

func (node *EchelonNode) Warnf(format string, args ...interface{}) {
	node.Logf(WarnLevel, format, args...)
}

func (node *EchelonNode) Errorf(format string, args ...interface{}) {
	node.Logf(ErrorLevel, format, args...)
}

func (node *EchelonNode) IsLevelEnabled(level LogLevel) bool {
	return node.config.Level >= level
}

func (node *EchelonNode) Logf(level LogLevel, format string, args ...interface{}) {
	node.AppendDescription(level, fmt.Sprintf(format, args...)+"\n")
}

func (node *EchelonNode) AppendDescription(level LogLevel, text string) {
	if !node.IsLevelEnabled(level) {
		return
	}
	node.lock.Lock()
	defer node.lock.Unlock()
	linesToAppend := strings.Split(text, "\n")
	if len(linesToAppend) == 0 {
		return
	}
	if len(node.description) == 0 {
		node.description = linesToAppend
		return
	}
	// append first new line to the last one
	node.description[len(node.description)-1] = node.description[len(node.description)-1] + linesToAppend[0]
	if len(linesToAppend) > 1 {
		node.description = append(node.description, linesToAppend[1:]...)
	}
}
