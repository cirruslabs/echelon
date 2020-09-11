package main

import (
	"fmt"
	"github.com/cirruslabs/echelon"
	"github.com/cirruslabs/echelon/renderers"
	"math/rand"
	"os"
	"sync/atomic"
	"time"
)

func main() {
	// renderer := renderers.NewSimpleRenderer(os.Stdout, nil)
	renderer := renderers.NewInteractiveRenderer(os.Stdout, nil)
	go renderer.StartDrawing()
	defer renderer.StopDrawing()
	log := echelon.NewLogger(echelon.InfoLevel, renderer)
	generateNode(log, 10)
	log.Finish(true)
}

//nolint:gochecknoglobals
var jobIDCounter uint64

func generateNode(log *echelon.Logger, magicConstant int) {
	jobID := atomic.AddUint64(&jobIDCounter, 1)
	scoped := log.Scoped(fmt.Sprintf("Job %d", jobID))
	for step := 0; step < magicConstant; step++ {
		//nolint:gosec,gomnd
		if rand.Intn(100) < magicConstant {
			generateNode(log, magicConstant-1)
		} else {
			childJobID := atomic.AddUint64(&jobIDCounter, 1)
			child := scoped.Scoped(fmt.Sprintf("Job %d", childJobID))
			subJobDuration := rand.Intn(magicConstant)
			for waitSecond := 0; waitSecond < subJobDuration; waitSecond++ {
				time.Sleep(time.Second)
				child.Infof("Doing very important jobs! Completed %d/100...", 100*(waitSecond+1)/subJobDuration)
			}
			child.Finish(true)
		}
	}
	scoped.Debugf("Finished after %d iterations", magicConstant)
	scoped.Finish(true)
}
