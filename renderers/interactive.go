package renderers

import (
	"bufio"
	"github.com/cirruslabs/echelon/logger"
	"github.com/cirruslabs/echelon/renderers/config"
	"github.com/cirruslabs/echelon/renderers/internal/node"
	"github.com/cirruslabs/echelon/terminal"
	"io"
	"sync"
	"time"
)

const resetAutoWrap = "\u001B[?7l"
const defaultFrameBufSize = 38400 // 80 by 120 of 4 bytes UTF-8 characters

type InteractiveRenderer struct {
	out               *bufio.Writer
	rootNode          *node.EchelonNode
	config            *config.InteractiveRendererConfig
	currentFrameLines []string
	drawLock          sync.Mutex
}

func NewInteractiveRenderer(out io.Writer, rendererConfig *config.InteractiveRendererConfig) *InteractiveRenderer {
	if rendererConfig == nil {
		rendererConfig = config.NewDefaultRenderingConfig()
	}
	return &InteractiveRenderer{
		out:      bufio.NewWriterSize(out, defaultFrameBufSize),
		rootNode: node.NewEchelonNode("root", rendererConfig),
		config:   rendererConfig,
	}
}

func findScopedNode(scopes []string, r *InteractiveRenderer) *node.EchelonNode {
	result := r.rootNode
	for _, scope := range scopes {
		result = result.FindOrCreateChild(scope)
	}
	return result
}

func (r *InteractiveRenderer) RenderScopeStarted(entry *logger.LogScopeStarted) {
	findScopedNode(entry.GetScopes(), r).Start()
}

func (r *InteractiveRenderer) RenderScopeFinished(entry *logger.LogScopeFinished) {
	n := findScopedNode(entry.GetScopes(), r)
	if entry.Success() {
		if n != r.rootNode {
			n.ClearAllChildren()
			n.ClearDescription()
		}
		// for some reason ✅ eats a space when displayed in terminal after it so let's add an extra space
		n.CompleteWithColor("✅ ", r.config.Colors.SuccessColor)
	} else {
		// for some reason ❌ eats a space when displayed in terminal after it so let's add an extra space
		n.CompleteWithColor("❌ ", r.config.Colors.FailureColor)
	}
}

func (r *InteractiveRenderer) RenderMessage(entry *logger.LogEntryMessage) {
	findScopedNode(entry.GetScopes(), r).AppendDescription(entry.GetMessage() + "\n")
}

func (r *InteractiveRenderer) StartDrawing() {
	// don't wrap lines since it breaks incremental redraws
	r.out.WriteString(resetAutoWrap)
	for !r.rootNode.HasCompleted() {
		r.DrawFrame()
		time.Sleep(r.config.RefreshRate)
	}
}

func (r *InteractiveRenderer) StopDrawing() {
	r.rootNode.Complete()
	// one last redraw
	r.DrawFrame()
}

func (r *InteractiveRenderer) DrawFrame() {
	r.drawLock.Lock()
	defer r.drawLock.Unlock()
	var newFrameLines []string
	for _, n := range r.rootNode.GetChildren() {
		newFrameLines = append(newFrameLines, n.Render()...)
	}
	terminal.CalculateIncrementalUpdate(r.out, r.currentFrameLines, newFrameLines)
	r.currentFrameLines = newFrameLines
}
