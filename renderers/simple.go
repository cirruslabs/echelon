package renderers

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/cirruslabs/echelon"
	"github.com/cirruslabs/echelon/renderers/internal/console"
	"github.com/cirruslabs/echelon/terminal"
	"github.com/cirruslabs/echelon/utils"
)

type SimpleRenderer struct {
	out        io.Writer
	colors     *terminal.ColorSchema
	startTimes map[string]time.Time
}

func NewSimpleRenderer(out io.Writer, colors *terminal.ColorSchema) *SimpleRenderer {
	if colors == nil {
		colors = terminal.DefaultColorSchema()
	}
	_ = console.PrepareTerminalEnvironment()
	return &SimpleRenderer{
		out:        out,
		colors:     colors,
		startTimes: make(map[string]time.Time),
	}
}

func (r SimpleRenderer) RenderScopeStarted(entry *echelon.LogScopeStarted) {
	scopes := entry.GetScopes()
	level := len(scopes)
	if level == 0 {
		return
	}
	timeKey := strings.Join(scopes, "/")
	if _, ok := r.startTimes[timeKey]; ok {
		// duplicate event
		return
	}
	r.startTimes[timeKey] = time.Now()
	lastScope := scopes[level-1]
	message := terminal.GetColoredText(r.colors.NeutralColor, fmt.Sprintf("Started %s", quotedIfNeeded(lastScope)))
	r.renderEntry(message)
}

func (r SimpleRenderer) RenderScopeFinished(entry *echelon.LogScopeFinished) {
	scopes := entry.GetScopes()
	level := len(scopes)
	if level == 0 {
		return
	}
	now := time.Now()
	startTime := now
	if t, ok := r.startTimes[strings.Join(scopes, "/")]; ok {
		startTime = t
	}
	duration := now.Sub(startTime)
	formatedDuration := utils.FormatDuration(duration, true)
	lastScope := scopes[level-1]

	switch entry.FinishType() {
	case echelon.FinishTypeSucceeded:
		message := fmt.Sprintf("%s succeeded in %s!", quotedIfNeeded(lastScope), formatedDuration)
		coloredMessage := terminal.GetColoredText(r.colors.SuccessColor, message)
		r.renderEntry(coloredMessage)
	case echelon.FinishTypeFailed:
		message := fmt.Sprintf("%s failed in %s!", quotedIfNeeded(lastScope), formatedDuration)
		coloredMessage := terminal.GetColoredText(r.colors.FailureColor, message)
		r.renderEntry(coloredMessage)
	case echelon.FinishTypeSkipped:
		message := fmt.Sprintf("%s skipped in %s!", quotedIfNeeded(lastScope), formatedDuration)
		coloredMessage := terminal.GetColoredText(r.colors.NeutralColor, message)
		r.renderEntry(coloredMessage)
	}
}

func (r SimpleRenderer) RenderMessage(entry *echelon.LogEntryMessage) {
	r.renderEntry(entry.GetMessage())
}

func (r SimpleRenderer) renderEntry(message string) {
	_, _ = r.out.Write([]byte(message + "\n"))
}

func (r SimpleRenderer) ScopeHasStarted(scopes []string) bool {
	level := len(scopes)
	if level == 0 {
		return true
	}
	timeKey := strings.Join(scopes, "/")
	_, result := r.startTimes[timeKey]
	return result
}

func quotedIfNeeded(s string) string {
	if strings.ContainsAny(s, "'\"") {
		return s
	}
	return "'" + s + "'"
}
