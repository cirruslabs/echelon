package node

import (
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

type EchelonNodeRenderingConfig struct {
	Level                          LogLevel
	ProgressIndicatorFrames        []string
	ProgressIndicatorCycleDuration time.Duration
}

func NewDefaultRenderingConfig() *EchelonNodeRenderingConfig {
	return &EchelonNodeRenderingConfig{
		Level: InfoLevel,
		ProgressIndicatorFrames: []string{
			"🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚", "🕛",
		},
		ProgressIndicatorCycleDuration: time.Second,
	}
}

func (config *EchelonNodeRenderingConfig) CurrentProgressIndicatorFrame() string {
	amountOfFrames := int64(len(config.ProgressIndicatorFrames))
	nanosPerFrame := int64(config.ProgressIndicatorCycleDuration) / amountOfFrames
	currentNanosTail := time.Now().UnixNano() % int64(config.ProgressIndicatorCycleDuration)
	frameIndex := currentNanosTail / nanosPerFrame
	if frameIndex < amountOfFrames {
		return config.ProgressIndicatorFrames[frameIndex]
	}
	return config.ProgressIndicatorFrames[0]
}
