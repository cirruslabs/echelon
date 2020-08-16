package echelon

type genericLogEntry struct {
	LogStarted  *LogScopeStarted
	LogFinished *LogScopeFinished
	LogEntry    *LogEntryMessage
}

type LogRendered interface {
	RenderScopeStarted(entry *LogScopeStarted)
	RenderScopeFinished(entry *LogScopeFinished)
	RenderMessage(entry *LogEntryMessage)
}

type Logger struct {
	level          LogLevel
	scopes         []string
	entriesChannel chan *genericLogEntry
}

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

func (logger *Logger) Finish(success bool) {
	logger.entriesChannel <- &genericLogEntry{
		LogFinished: NewLogScopeFinished(success, logger.scopes...),
	}
}

func (logger *Logger) IsLogLevelEnabled(level LogLevel) bool {
	return level <= logger.level
}
