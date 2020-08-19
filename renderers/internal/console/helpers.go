// +build !windows

package console

func PrepareTerminalEnvironment() error {
	// no need on unix
	return nil
}
