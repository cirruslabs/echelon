package main

import (
	"fmt"
	"github.com/cirruslabs/echelon/console"
	"github.com/cirruslabs/echelon/node"
	"math/rand"
	"os"
	"sync/atomic"
	"time"
)

func main() {
	config := node.NewDefaultRenderingConfig()
	console.NewConsole(os.Stdout, []*node.EchelonNode{generateNode(config, 10)}).StartDrawing()
}

var jobIdCounter uint64

func generateNode(config *node.EchelonNodeConfig, magicConstant int) *node.EchelonNode {
	jobId := atomic.AddUint64(&jobIdCounter, 1)
	result := node.StartNewEchelonNode(fmt.Sprintf("Job %d", jobId), config)
	go func() {
		for step := 0; step < magicConstant; step++ {
			if rand.Intn(100) < magicConstant {
				child := generateNode(config, magicConstant-1)
				result.AddNewChild(child)
				child.WaitCompletion()
			} else {
				childJobId := atomic.AddUint64(&jobIdCounter, 1)
				child := result.StartNewChild(fmt.Sprintf("Job %d", childJobId))
				subJobDuration := rand.Intn(magicConstant)
				for waitSecond := 0; waitSecond < subJobDuration; waitSecond++ {
					time.Sleep(time.Second)
					child.Infof("Doing very important jobs! Completed %d/100...", 100*(waitSecond+1)/subJobDuration)
				}
				child.ClearDescription()
				child.SetTitleColor(node.GREEN_COLOR)
				child.SetStatus("👍")
				child.Complete()
			}
		}
		result.ClearAllChildren()
		result.SetTitleColor(node.GREEN_COLOR)
		result.SetStatus("👍")
		result.Complete()
	}()
	return result
}
