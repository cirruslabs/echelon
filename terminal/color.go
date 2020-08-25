package terminal

import "fmt"

type ColorSchema struct {
	SuccessColor int
	FailureColor int
	NeutralColor int
}

// Reset ANSI sequence.
const ResetSequence = "\033[0m"

const (
	BlackColor = iota
	RedColor
	GreenColor
	YellowColor
	BlueColor
	MagentaColor
	CyanColor
	WhiteColor
)

func DefaultColorSchema() *ColorSchema {
	return &ColorSchema{
		SuccessColor: GreenColor,
		FailureColor: RedColor,
		NeutralColor: YellowColor,
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
