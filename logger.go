package echelon

type genericLogEntry struct {
	LogStarted  *LogScopeStarted
	LogFinished *LogScopeFinished
	LogEntry    *LogEntryMessage
	Annotation  *Annotation
}

type LogRendered interface {
	RenderScopeStarted(entry *LogScopeStarted)
	RenderScopeFinished(entry *LogScopeFinished)
	RenderMessage(entry *LogEntryMessage)
	RenderAnnotation(entry *Annotation)
}

type Logger struct {
	level          LogLevel
	scopes         []string
	entriesChannel chan *genericLogEntry
}

type FinishType int

const (
	FinishTypeSucceeded FinishType = iota
	FinishTypeFailed
	FinishTypeSkipped
)

func NewLogger(level LogLevel, renderer LogRendered) *Logger {
	logger := &Logger{
		level:          level,
		entriesChannel: make(chan *genericLogEntry),
	}
	go logger.streamEntries(renderer)
	return logger
}

func (logger *Logger) Scoped(scope string) *Logger {
	result := &Logger{
		level:          logger.level,
		scopes:         append(logger.scopes, scope),
		entriesChannel: logger.entriesChannel,
	}
	result.entriesChannel <- &genericLogEntry{
		LogStarted: NewLogScopeStarted(result.scopes...),
	}
	return result
}

func (logger *Logger) streamEntries(renderer LogRendered) {
	for {
		entry := <-logger.entriesChannel
		if entry.LogStarted != nil {
			renderer.RenderScopeStarted(entry.LogStarted)
		}
		if entry.LogFinished != nil {
			renderer.RenderScopeFinished(entry.LogFinished)
		}
		if entry.LogEntry != nil {
			renderer.RenderMessage(entry.LogEntry)
		}
		if entry.Annotation != nil {
			renderer.RenderAnnotation(entry.Annotation)
		}
	}
}

func (logger *Logger) Tracef(format string, args ...interface{}) {
	logger.Logf(TraceLevel, format, args...)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.Logf(DebugLevel, format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.Logf(InfoLevel, format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.Logf(WarnLevel, format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.Logf(ErrorLevel, format, args...)
}

func (logger *Logger) Logf(level LogLevel, format string, args ...interface{}) {
	if logger.IsLogLevelEnabled(level) {
		logger.entriesChannel <- &genericLogEntry{
			LogEntry: NewLogEntryMessage(logger.scopes, level, format, args...),
		}
	}
}

type AnnotationLevel int

const (
	AnnotationLevelNotice AnnotationLevel = iota
	AnnotationLevelWarning
	AnnotationLevelError
)

type Annotation struct {
	Level     AnnotationLevel
	File      string
	LineStart int64
	LineEnd   int64
	Title     string
	Message   string
}

func (logger *Logger) Annotation(annotation *Annotation) {
	logger.entriesChannel <- &genericLogEntry{
		Annotation: annotation,
	}
}

func (logger *Logger) Finish(success bool) {
	var finishType FinishType

	if success {
		finishType = FinishTypeSucceeded
	} else {
		finishType = FinishTypeFailed
	}

	logger.FinishWithType(finishType)
}

func (logger *Logger) FinishWithType(finishType FinishType) {
	logger.entriesChannel <- &genericLogEntry{
		LogFinished: NewLogScopeFinished(finishType, logger.scopes...),
	}
}

func (logger *Logger) IsLogLevelEnabled(level LogLevel) bool {
	return level <= logger.level
}
