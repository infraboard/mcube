package file

import (
	"os"

	"golang.org/x/sys/windows"
)

// RedirectStandardError causes all standard error output to be directed to the
// given file.
func RedirectStandardError(toFile *os.File) error {
	return windows.SetStdHandle(windows.STD_ERROR_HANDLE, windows.Handle(toFile.Fd()))
}
