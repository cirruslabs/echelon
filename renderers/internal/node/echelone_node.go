package node

import (
	"fmt"
	"github.com/cirruslabs/echelon/renderers/config"
	"github.com/cirruslabs/echelon/terminal"
	"math"
	"strings"
	"sync"
	"time"
)

type EchelonNode struct {
	lock                    sync.RWMutex
	done                    sync.WaitGroup
	status                  string
	title                   string
	titleColor              int
	description             []string
	visibleDescriptionLines int
	config                  *config.InteractiveRendererConfig
	startTime               time.Time
	endTime                 time.Time
	children                []*EchelonNode
}

func StartNewEchelonNode(title string, config *config.InteractiveRendererConfig) *EchelonNode {
	result := NewEchelonNode(title, config)
	result.Start()
	return result
}

func NewEchelonNode(title string, config *config.InteractiveRendererConfig) *EchelonNode {
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

func (node *EchelonNode) GetChildren() []*EchelonNode {
	node.lock.RLock()
	defer node.lock.RUnlock()
	return node.children
}

func (node *EchelonNode) UpdateTitle(text string) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.title = text
}

func (node *EchelonNode) UpdateConfig(config *config.InteractiveRendererConfig) {
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
		coloredTitle = terminal.GetColoredText(node.titleColor, node.title)
	}
	duration := formatDuration(node.ExecutionDuration(), len(node.children) == 0)
	return fmt.Sprintf("%s %s %s", prefix, coloredTitle, duration)
}

func formatDuration(duration time.Duration, showDecimals bool) string {
	if duration < 10*time.Second && showDecimals {
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
	return !node.endTime.IsZero()
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
	if node.startTime.IsZero() {
		node.startTime = node.endTime
	}
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
	if node.startTime.IsZero() {
		node.startTime = node.endTime
	}
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

func (node *EchelonNode) AppendDescription(text string) {
	if node.HasCompleted() {
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
