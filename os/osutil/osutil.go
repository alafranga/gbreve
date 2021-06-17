package osutil

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// IsTerminal should be commented
func IsTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
