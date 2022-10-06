package echelon

import (
	"fmt"
	"time"
)

type LogLevel uint32

const (
	ErrorLevel LogLevel = iota
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type LogScopeStarted struct {
	scopes []string
	time   time.Time
}

func NewLogScopeStarted(scopes ...string) *LogScopeStarted {
	return &LogScopeStarted{
		scopes: scopes,
		time:   time.Now(),
	}
}

func (entry *LogScopeStarted) GetScopes() []string {
	return entry.scopes
}

type LogScopeFinished struct {
	scopes     []string
	finishType FinishType
}

func NewLogScopeFinished(finishType FinishType, scopes ...string) *LogScopeFinished {
	return &LogScopeFinished{
		scopes:     scopes,
		finishType: finishType,
	}
}

func (entry *LogScopeFinished) FinishType() FinishType {
	return entry.finishType
}

func (entry *LogScopeFinished) GetScopes() []string {
	return entry.scopes
}

type LogEntryMessage struct {
	Level   LogLevel
	scopes  []string
	message string
	raw     bool
}

func NewLogEntryMessage(scopes []string, level LogLevel, format string, arguments ...interface{}) *LogEntryMessage {
	return &LogEntryMessage{
		Level:   level,
		scopes:  scopes,
		message: fmt.Sprintf(format, arguments...),
	}
}

func (entry *LogEntryMessage) GetMessage() string {
	if entry.raw {
		return entry.message
	}
	return entry.message + "\n"
}

func (entry *LogEntryMessage) GetScopes() []string {
	return entry.scopes
}
