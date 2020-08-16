package utils

import (
	"fmt"
	"math"
	"time"
)

func FormatDuration(duration time.Duration, showDecimals bool) string {
	if duration < 10*time.Second && showDecimals {
		return fmt.Sprintf("%.1fs", float64(duration.Milliseconds())/1000)
	}
	if duration < time.Minute {
		return fmt.Sprintf("%ds", int(math.Floor(duration.Seconds()))%60)
	}
	if duration < time.Hour {
		return fmt.Sprintf("%02d:%02d", int(math.Floor(duration.Minutes()))%60, int(math.Floor(duration.Seconds()))%60)
	}
	return fmt.Sprintf(
		"%02d:%02d:%02d",
		int(math.Floor(duration.Hours())),
		int(math.Floor(duration.Minutes()))%60,
		int(math.Floor(duration.Seconds()))%60,
	)
}
