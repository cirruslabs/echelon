package terminal

import "fmt"

type ColorSchema struct {
	SuccessColor int
	FailureColor int
	NeutralColor int
}

// Reset ANSI sequence
const ResetSequence = "\033[0m"

const (
	BLACK_COLOR = iota
	RED_COLOR
	GREEN_COLOR
	YELLOW_COLOR
	BLUE_COLOR
	MAGENTA_COLOR
	CYAN_COLOR
	WHITE_COLOR
)

func DefaultColorSchema() *ColorSchema {
	return &ColorSchema{
		SuccessColor: GREEN_COLOR,
		FailureColor: RED_COLOR,
		NeutralColor: YELLOW_COLOR,
	}
}

func GetColoredText(color int, text string) string {
	return fmt.Sprintf("%s%s%s", GetColorSequence(color), text, ResetSequence)
}

func GetColorSequence(code int) string {
	if code < 0 {
		return ResetSequence
	}
	return fmt.Sprintf("\033[3%dm", code)
}
