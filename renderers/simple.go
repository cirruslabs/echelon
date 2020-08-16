package renderers

import (
	"fmt"
	"github.com/cirruslabs/echelon"
	"github.com/cirruslabs/echelon/terminal"
	"github.com/cirruslabs/echelon/utils"
	"io"
	"strings"
	"time"
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
	r.startTimes[strings.Join(scopes, "/")] = time.Now()
	lastScope := scopes[level-1]
	message := terminal.GetColoredText(r.colors.NeutralColor, fmt.Sprintf("Started '%s'", lastScope))
	r.renderEntryWithIndention(level-1, message)
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
	if entry.Success() {
		message := fmt.Sprintf("'%s' succeded in %s!", lastScope, formatedDuration)
		coloredMessage := terminal.GetColoredText(r.colors.SuccessColor, message)
		r.renderEntryWithIndention(level, coloredMessage)
	} else {
		message := fmt.Sprintf("'%s' failed in %s!", lastScope, formatedDuration)
		coloredMessage := terminal.GetColoredText(r.colors.NeutralColor, message)
		r.renderEntryWithIndention(level, coloredMessage)
	}
}

func (r SimpleRenderer) RenderMessage(entry *echelon.LogEntryMessage) {
	r.renderEntryWithIndention(len(entry.GetScopes()), entry.GetMessage())
}

func (r SimpleRenderer) renderEntryWithIndention(level int, message string) {
	if level <= 0 {
		_, _ = r.out.Write([]byte(message + "\n"))
	} else {
		prefix := strings.Repeat("  ", level)
		_, _ = r.out.Write(append([]byte(prefix), []byte(message+"\n")...))
	}
}
