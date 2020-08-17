package config

import (
	"github.com/cirruslabs/echelon/terminal"
	"time"
)

type InteractiveRendererConfig struct {
	Colors                         *terminal.ColorSchema
	RefreshRate                    time.Duration
	ProgressIndicatorFrames        []string
	ProgressIndicatorCycleDuration time.Duration
	SuccessStatus                  string
	FailureStatus                  string
}

func NewDefaultRenderingConfig() *InteractiveRendererConfig {
	return &InteractiveRendererConfig{
		Colors:      terminal.DefaultColorSchema(),
		RefreshRate: 200 * time.Microsecond,
		ProgressIndicatorFrames: []string{
			"🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚", "🕛",
		},
		ProgressIndicatorCycleDuration: time.Second,
		SuccessStatus:                  "✅",
		FailureStatus:                  "❌",
	}
}

func (config *InteractiveRendererConfig) CurrentProgressIndicatorFrame() string {
	amountOfFrames := int64(len(config.ProgressIndicatorFrames))
	nanosPerFrame := int64(config.ProgressIndicatorCycleDuration) / amountOfFrames
	currentNanosTail := time.Now().UnixNano() % int64(config.ProgressIndicatorCycleDuration)
	frameIndex := currentNanosTail / nanosPerFrame
	if frameIndex < amountOfFrames {
		return config.ProgressIndicatorFrames[frameIndex]
	}
	return config.ProgressIndicatorFrames[0]
}
