package utils

import (
	"fmt"
	"math"
	"time"
)

const (
	millisInSecond  = 1000
	secondsInMinute = 60
	minutesInHour   = 60
)

func FormatDuration(duration time.Duration, showDecimals bool) string {
	if duration < 10*time.Second && showDecimals {
		return fmt.Sprintf("%.1fs", float64(duration.Milliseconds())/millisInSecond)
	}
	seconds := int(math.Floor(duration.Seconds())) % secondsInMinute
	if duration < time.Minute {
		return fmt.Sprintf("%ds", seconds)
	}
	minutes := int(math.Floor(duration.Minutes())) % minutesInHour
	if duration < time.Hour {
		return fmt.Sprintf("%02d:%02d", minutes, seconds)
	}
	hours := int(math.Floor(duration.Hours()))
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
