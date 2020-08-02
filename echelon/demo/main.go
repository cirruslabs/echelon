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
	console.NewConsole(os.Stdout, generateNode(10)).StartDrawing()
}

var jobIdCounter uint64

func generateNode(magicConstant int) *node.EchelonNode {
	jobId := atomic.AddUint64(&jobIdCounter, 1)
	result := node.StartNewEchelonNode(fmt.Sprintf("Job %d", jobId))
	go func() {
		for step := 0; step < magicConstant; step++ {
			if rand.Intn(100) < magicConstant {
				child := generateNode(magicConstant - 1)
				result.AddNewChild(child)
				child.Wait()
			} else {
				childJobId := atomic.AddUint64(&jobIdCounter, 1)
				child := node.StartNewEchelonNode(fmt.Sprintf("Job %d", childJobId))
				result.AddNewChild(child)
				time.Sleep(time.Duration(rand.Intn(magicConstant)) * time.Second)
				child.Complete()
			}
		}
		result.Complete()
	}()
	return result
}
