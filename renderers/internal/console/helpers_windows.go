package console

import (
	"golang.org/x/sys/windows"
	"os"
)

func PrepareTerminalEnvironment() error {
	// enable handling ASCII codes
	err := addConsoleMode(windows.Stdout, windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	if err != nil {
		return err
	}
	return addConsoleMode(windows.Stderr, windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}

func addConsoleMode(handle windows.Handle, flags uint32) error {
	var mode uint32

	err := windows.GetConsoleMode(handle, &mode)
	if err != nil {
		return err
	}
	return windows.SetConsoleMode(handle, mode|flags)
}

func TerminalHeight(file *os.File) int {
	// todo: figure out how to find out console height on Windows
	return -1
}
