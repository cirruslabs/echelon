package renderers

import (
	"fmt"
	"github.com/cirruslabs/echelon/logger"
	"github.com/cirruslabs/echelon/terminal"
	"io"
	"strings"
)

type SimpleRenderer struct {
	out    io.Writer
	colors *terminal.ColorSchema
}

func NewSimpleRenderer(out io.Writer, colors *terminal.ColorSchema) *SimpleRenderer {
	if colors == nil {
		colors = terminal.DefaultColorSchema()
	}
	return &SimpleRenderer{
		out:    out,
		colors: colors,
	}
}

func (r SimpleRenderer) RenderScopeStarted(entry *logger.LogScopeStarted) {
	scopes := entry.GetScopes()
	level := len(scopes)
	if level == 0 {
		return
	}
	lastScope := scopes[level-1]
	message := terminal.GetColoredText(r.colors.NeutralColor, fmt.Sprintf("Started '%s'", lastScope))
	r.renderEntryWithIndention(level-1, message)
}

func (r SimpleRenderer) RenderScopeFinished(entry *logger.LogScopeFinished) {
	scopes := entry.GetScopes()
	level := len(scopes)
	if level == 0 {
		return
	}
	lastScope := scopes[level-1]
	if entry.Success() {
		message := terminal.GetColoredText(r.colors.SuccessColor, fmt.Sprintf("'%s' succeded!", lastScope))
		r.renderEntryWithIndention(level, message)
	} else {
		message := terminal.GetColoredText(r.colors.NeutralColor, fmt.Sprintf("'%s' has failed!", lastScope))
		r.renderEntryWithIndention(level, message)
	}
}

func (r SimpleRenderer) RenderMessage(entry *logger.LogEntryMessage) {
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
