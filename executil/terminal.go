package executil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var PasswordPrompt = `Enter password: `
var PasswordVerifyPrompt = `Verify password: `
var PromptWriter io.Writer = os.Stdout

// Read a password from standard input, disabling echo.
func ReadPassword() string {
	pw, _ := PromptPassword(PromptWriter, int(syscall.Stdin), false)
	return pw
}

// Read a password from standard input, disabling echo and prompting twice.  The
// second return argument is falseif the two passwords do not match.
func ReadAndVerifyPassword() (string, bool) {
	return PromptPassword(PromptWriter, int(syscall.Stdin), true)
}

// Generic password prompt that takes the input file descriptor and verify flag as options.
func PromptPassword(writer io.Writer, fd int, verify bool) (string, bool) {
	fmt.Fprint(writer, PasswordPrompt)

	if pw1, err := terminal.ReadPassword(fd); err == nil {
		fmt.Fprint(writer, "\n")

		if verify {
			fmt.Fprint(writer, PasswordVerifyPrompt)

			if pw2, err := terminal.ReadPassword(fd); err == nil {
				fmt.Fprint(writer, "\n")

				if bytes.Equal(pw1, pw2) {
					return string(pw1), true
				} else {
					return string(pw1), false
				}
			}
		} else {
			return string(pw1), true
		}
	}

	return ``, false
}
