// +build !windows

package console

import (
	"golang.org/x/sys/unix"
	"os"
)

func PrepareTerminalEnvironment() error {
	// no need on unix
	return nil
}

func TerminalHeight(file *os.File) int {
	ws, err := unix.IoctlGetWinsize(int(file.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return -1
	}

	return int(ws.Row)
}
