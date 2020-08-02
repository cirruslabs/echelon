package node

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type EchelonNode struct {
	lock                sync.RWMutex
	done                sync.WaitGroup
	title               string
	description         []string
	maxDescriptionLines int
	startTime           time.Time
	endTime             time.Time
	children            []*EchelonNode
}

func StartNewEchelonNode(title string) *EchelonNode {
	result := &EchelonNode{
		title:               title,
		description:         make([]string, 0),
		maxDescriptionLines: 5,
		startTime:           time.Now(),
		endTime:             time.Unix(0, 0),
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

func (node *EchelonNode) AppendDescription(text string) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.description = append(node.description, text)
	linesTotal := len(node.description)
	if linesTotal > node.maxDescriptionLines {
		node.description = node.description[(linesTotal - node.maxDescriptionLines):]
	}
}

func (node *EchelonNode) Draw() []string {
	node.lock.RLock()
	defer node.lock.RUnlock()
	result := []string{
		fmt.Sprintf("%s %s", node.title, formatDuration(node.ExecutionDuration())),
	}
	if len(node.children) > 0 {
		for _, child := range node.children {
			for _, childDescriptionLine := range child.Draw() {
				result = append(result, "  "+childDescriptionLine)
			}
		}
	} else {
		for _, descriptionLine := range node.description {
			result = append(result, "  "+descriptionLine)
		}
	}

	return result
}

func formatDuration(duration time.Duration) string {
	if duration < 10*time.Second {
		return fmt.Sprintf("%.1fs", float64(duration.Milliseconds())/1000)
	}
	if duration < time.Minute {
		return fmt.Sprintf("%ds", int(math.Floor(duration.Seconds())))
	}
	if duration < time.Hour {
		return fmt.Sprintf("%02d:%02d", int(math.Floor(duration.Minutes())), int(math.Floor(duration.Seconds())))
	}
	return fmt.Sprintf(
		"%02d:%02d:%02d",
		int(math.Floor(duration.Hours())),
		int(math.Floor(duration.Minutes())),
		int(math.Floor(duration.Seconds())),
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

func (node *EchelonNode) IsRunning() bool {
	node.lock.RLock()
	defer node.lock.RUnlock()
	return node.endTime.Before(node.startTime)
}

func (node *EchelonNode) AddNewChild(child *EchelonNode) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.children = append(node.children, child)
}

func (node *EchelonNode) Complete() {
	node.lock.Lock()
	node.endTime = time.Now()
	node.lock.Unlock()
	node.done.Done()
}

func (node *EchelonNode) Wait() {
	node.done.Wait()
}
