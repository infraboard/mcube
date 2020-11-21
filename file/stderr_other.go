// +build !windows

package file

import (
	"os"

	"golang.org/x/sys/unix"
)

// RedirectStandardError causes all standard error output to be directed to the
// given file.
func RedirectStandardError(toFile *os.File) error {
	return unix.Dup2(int(toFile.Fd()), 2)
}
