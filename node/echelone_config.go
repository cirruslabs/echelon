package node

import (
	"time"
)

type EchelonNodeRenderingConfig struct {
	ProgressIndicatorFrames        []string
	ProgressIndicatorCycleDuration time.Duration
	MaxDescriptionLines            int
}

func NewDefaultRenderingConfig() *EchelonNodeRenderingConfig {
	return &EchelonNodeRenderingConfig{
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
