package renderers

import "github.com/cirruslabs/echelon"

type StubRenderer struct{}

func (*StubRenderer) RenderScopeStarted(entry *echelon.LogScopeStarted) {}

func (*StubRenderer) RenderScopeFinished(entry *echelon.LogScopeFinished) {}

func (*StubRenderer) RenderMessage(entry *echelon.LogEntryMessage) {}

func (*StubRenderer) RenderAnnotation(entry *echelon.Annotation) {}
