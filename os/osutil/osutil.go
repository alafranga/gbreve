package osutil

import (
	"os"

	terminal "golang.org/x/term"
)

// IsTerminal should be commented
func IsTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
