package node

import (
	"fmt"
	"github.com/cirruslabs/echelon/renderers/config"
	"github.com/cirruslabs/echelon/terminal"
	"github.com/cirruslabs/echelon/utils"
	"golang.org/x/text/width"
	"strings"
	"sync"
	"time"
)

const defaultVisibleLines = 5

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
		titleColor:              config.Colors.NeutralColor,
		description:             make([]string, 0),
		visibleDescriptionLines: defaultVisibleLines,
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
	title := node.fancyTitle()
	tail := node.renderChildren()
	node.lock.RLock()
	defer node.lock.RUnlock()
	if len(node.description) > node.visibleDescriptionLines && node.visibleDescriptionLines >= 0 {
		tail = append(tail, "...")
		tail = append(tail, node.description[(len(node.description)-node.visibleDescriptionLines):]...)
	} else {
		tail = append(tail, node.description...)
	}
	indent := "  " // two spaces by default
	props, _ := width.LookupString(title)
	if props.Kind() == width.EastAsianWide || props.Kind() == width.EastAsianFullwidth {
		indent = "   " // three spaces since title start with a wide emoji
	}
	result := []string{title}
	for _, descriptionLine := range tail {
		result = append(result, indent+descriptionLine)
	}

	return result
}

func (node *EchelonNode) renderChildren() []string {
	node.lock.RLock()
	defer node.lock.RUnlock()
	var result []string
	for _, child := range node.children {
		result = append(result, child.Render()...)
	}
	return result
}

func (node *EchelonNode) fancyTitle() string {
	duration := utils.FormatDuration(node.ExecutionDuration(), len(node.children) == 0)
	isRunning := node.IsRunning()

	node.lock.RLock()
	defer node.lock.RUnlock()
	prefix := node.status
	if isRunning {
		prefix = node.config.CurrentProgressIndicatorFrame()
	}
	coloredTitle := node.title
	if node.titleColor >= 0 {
		coloredTitle = terminal.GetColoredText(node.titleColor, node.title)
	}
	return fmt.Sprintf("%s %s %s", prefix, coloredTitle, duration)
}

func (node *EchelonNode) ExecutionDuration() time.Duration {
	node.lock.RLock()
	defer node.lock.RUnlock()
	if !node.startTime.IsZero() && node.endTime.IsZero() {
		return time.Since(node.startTime)
	}
	return node.endTime.Sub(node.startTime)
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
