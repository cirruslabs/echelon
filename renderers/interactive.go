package renderers

import (
	"bufio"
	"github.com/cirruslabs/echelon"
	"github.com/cirruslabs/echelon/renderers/config"
	"github.com/cirruslabs/echelon/renderers/internal/console"
	"github.com/cirruslabs/echelon/renderers/internal/node"
	"github.com/cirruslabs/echelon/terminal"
	"os"
	"sync"
	"time"
)

const disableAutoWrap = "\u001B[?7l"
const enableAutoWrap = "\u001B[?7h"
const defaultFrameBufSize = 38400 // 80 by 120 of 4 bytes UTF-8 characters

type InteractiveRenderer struct {
	out               *bufio.Writer
	rootNode          *node.EchelonNode
	config            *config.InteractiveRendererConfig
	currentFrameLines []string
	drawLock          sync.Mutex
	terminalHeight    int

	StubRenderer
}

func NewInteractiveRenderer(out *os.File, rendererConfig *config.InteractiveRendererConfig) *InteractiveRenderer {
	if rendererConfig == nil {
		rendererConfig = config.NewDefaultRenderingConfig()
	}
	return &InteractiveRenderer{
		out:            bufio.NewWriterSize(out, defaultFrameBufSize),
		rootNode:       node.NewEchelonNode("root", rendererConfig),
		config:         rendererConfig,
		terminalHeight: console.TerminalHeight(out),
	}
}

func findScopedNode(scopes []string, r *InteractiveRenderer) *node.EchelonNode {
	result := r.rootNode
	for _, scope := range scopes {
		result = result.FindOrCreateChild(scope)
	}
	return result
}

func (r *InteractiveRenderer) RenderScopeStarted(entry *echelon.LogScopeStarted) {
	findScopedNode(entry.GetScopes(), r).Start()
}

func (r *InteractiveRenderer) RenderScopeFinished(entry *echelon.LogScopeFinished) {
	n := findScopedNode(entry.GetScopes(), r)

	switch entry.FinishType() {
	case echelon.FinishTypeSucceeded:
		if n != r.rootNode {
			n.ClearAllChildren()
			n.ClearDescription()
		}
		n.CompleteWithColor(r.config.SuccessStatus, r.config.Colors.SuccessColor)
	case echelon.FinishTypeFailed:
		n.SetVisibleDescriptionLines(r.config.DescriptionLinesWhenFailed)
		n.CompleteWithColor(r.config.FailureStatus, r.config.Colors.FailureColor)
	case echelon.FinishTypeSkipped:
		if n != r.rootNode {
			n.ClearAllChildren()
			n.ClearDescription()
		}
		n.CompleteWithColor(r.config.SkippedStatus, r.config.Colors.NeutralColor)
	}
}

func (r *InteractiveRenderer) RenderMessage(entry *echelon.LogEntryMessage) {
	findScopedNode(entry.GetScopes(), r).AppendDescription(entry.GetMessage())
}

func (r *InteractiveRenderer) StartDrawing() {
	_ = console.PrepareTerminalEnvironment()
	// don't wrap lines since it breaks incremental redraws
	_, _ = r.out.WriteString(disableAutoWrap)
	for !r.rootNode.HasCompleted() {
		r.DrawFrame()
		time.Sleep(r.config.RefreshRate)
	}
}

func (r *InteractiveRenderer) StopDrawing() {
	r.rootNode.Complete()
	// one last redraw
	r.DrawFrame()
	// don't leave autowrap disabled in the terminal
	_, _ = r.out.WriteString(enableAutoWrap)
	_ = r.out.Flush()
}

func (r *InteractiveRenderer) DrawFrame() {
	r.drawLock.Lock()
	defer r.drawLock.Unlock()
	var newFrameLines []string
	for _, n := range r.rootNode.GetChildren() {
		newFrameLines = append(newFrameLines, n.Render()...)
	}
	if r.terminalHeight > 0 {
		terminal.CalculateIncrementalUpdateMaxLines(r.out, r.currentFrameLines, newFrameLines, r.terminalHeight)
	} else {
		terminal.CalculateIncrementalUpdate(r.out, r.currentFrameLines, newFrameLines)
	}
	r.currentFrameLines = newFrameLines
}
