package config

import (
	"github.com/cirruslabs/echelon/terminal"
	"runtime"
	"time"
)

type InteractiveRendererConfig struct {
	Colors                         *terminal.ColorSchema
	RefreshRate                    time.Duration
	ProgressIndicatorFrames        []string
	ProgressIndicatorCycleDuration time.Duration
	SuccessStatus                  string
	FailureStatus                  string
	DescriptionLinesWhenFailed     int
}

func NewDefaultRenderingConfig() *InteractiveRendererConfig {
	if runtime.GOOS == "windows" {
		return NewDefaultSymbolsOnlyRenderingConfig()
	}
	return NewDefaultEmojiRenderingConfig()
}

func NewDefaultEmojiRenderingConfig() *InteractiveRendererConfig {
	//nolint:gomnd
	return &InteractiveRendererConfig{
		Colors:      terminal.DefaultColorSchema(),
		RefreshRate: 200 * time.Microsecond,
		ProgressIndicatorFrames: []string{
			"🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚", "🕛",
		},
		ProgressIndicatorCycleDuration: time.Second,
		SuccessStatus:                  "✅",
		FailureStatus:                  "❌",
		DescriptionLinesWhenFailed:     100,
	}
}

func NewDefaultSymbolsOnlyRenderingConfig() *InteractiveRendererConfig {
	//nolint:gomnd
	return &InteractiveRendererConfig{
		Colors:      terminal.DefaultColorSchema(),
		RefreshRate: 250 * time.Microsecond,
		ProgressIndicatorFrames: []string{
			"\\", "|", "/", "-",
		},
		ProgressIndicatorCycleDuration: time.Second,
		SuccessStatus:                  "+",
		FailureStatus:                  "-",
		DescriptionLinesWhenFailed:     100,
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
