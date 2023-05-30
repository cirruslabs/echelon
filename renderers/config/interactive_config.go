package config

import (
	"github.com/cirruslabs/echelon/terminal"
	"runtime"
	"time"
)

const defaultVisibleLines = 5

type InteractiveRendererConfig struct {
	Colors                         *terminal.ColorSchema
	RefreshRate                    time.Duration
	ProgressIndicatorFrames        []string
	ProgressIndicatorCycleDuration time.Duration
	SuccessStatus                  string
	FailureStatus                  string
	SkippedStatus                  string
	DescriptionLinesWhenFailed     int
	DescriptionLinesWhenSkipped    int
	VisibleDescriptionLines        int
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
		SkippedStatus:                  "⏩",
		DescriptionLinesWhenFailed:     100,
		DescriptionLinesWhenSkipped:    0,
		VisibleDescriptionLines:        defaultVisibleLines,
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
		SkippedStatus:                  "!",
		DescriptionLinesWhenFailed:     100,
		DescriptionLinesWhenSkipped:    0,
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
