package main

import (
	"fmt"
	"github.com/cirruslabs/echelon/logger"
	"github.com/cirruslabs/echelon/renderers"
	"math/rand"
	"os"
	"sync/atomic"
	"time"
)

func main() {
	//renderer := renderers.NewSimpleRenderer(os.Stdout, nil)
	renderer := renderers.NewInteractiveRenderer(os.Stdout, nil)
	go renderer.StartDrawing()
	defer renderer.StopDrawing()
	log := logger.NewLogger(renderer)
	generateNode(log, 10)
	log.Finish(true)
}

var jobIdCounter uint64

func generateNode(log *logger.Logger, magicConstant int) {
	jobId := atomic.AddUint64(&jobIdCounter, 1)
	scoped := log.Scoped(fmt.Sprintf("Job %d", jobId))
	for step := 0; step < magicConstant; step++ {
		if rand.Intn(100) < magicConstant {
			generateNode(log, magicConstant-1)
		} else {
			childJobId := atomic.AddUint64(&jobIdCounter, 1)
			child := scoped.Scoped(fmt.Sprintf("Job %d", childJobId))
			subJobDuration := rand.Intn(magicConstant)
			for waitSecond := 0; waitSecond < subJobDuration; waitSecond++ {
				time.Sleep(time.Second)
				child.Infof("Doing very important jobs! Completed %d/100...", 100*(waitSecond+1)/subJobDuration)
			}
			child.Finish(true)
		}
	}
	scoped.Finish(true)
}
