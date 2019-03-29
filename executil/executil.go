package executil

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ghetzel/go-stockutil/fileutil"
	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

// Locates the first path containing the given command. The directories listed
// in the environment variable "PATH" will be checked, in order.  If additional
// directories are specified in the path variadic argument, they will be checked
// first.  If the command is not in any path, an empty string will be returned.
func Which(cmdname string, path ...string) string {
	if found := WhichAll(cmdname, path...); len(found) > 0 {
		return found[0]
	} else {
		return ``
	}
}

// Locates the all paths containing the given command. The directories listed
// in the environment variable "PATH" will be checked, in order.  If additional
// directories are specified in the path variadic argument, they will be checked
// first.  If the command is not in any path, an empty slice will be returned.
func WhichAll(cmdname string, path ...string) []string {
	dirs := append(path, strings.Split(os.Getenv(`PATH`), `:`)...)
	found := make([]string, 0)

	if fileutil.IsNonemptyExecutableFile(cmdname) {
		found = append(found, cmdname)
	}

	for _, dir := range dirs {
		candidate := filepath.Join(dir, cmdname)

		if len(strings.TrimSpace(dir)) == 0 {
			continue
		} else if !fileutil.DirExists(dir) {
			continue
		} else if fileutil.IsNonemptyExecutableFile(candidate) {
			found = append(found, candidate)
		}
	}

	return found
}

// Take an *exec.Cmd or []string and return a shell-executable command line string.
func Join(in interface{}) string {
	var args []string

	if cmd, ok := in.(*exec.Cmd); ok {
		args = cmd.Args
	} else if typeutil.IsArray(in) {
		args = sliceutil.Stringify(in)
	} else {
		return ``
	}

	for i, arg := range args {
		if strings.Contains(arg, ` `) {
			args[i] = stringutil.Wrap(arg, `'`, `'`)
		}
	}

	return strings.Join(args, ` `)
}

// Uses environment variables and other configurations to attempt to locate the
// path to the user's shell.
func FindShell() string {
	shells := []string{os.Getenv(`SHELL`)}
	shells = append(shells, Which(`sh`))

	for _, shell := range shells {
		if shell != `` {
			return shell
		}
	}

	return ``
}

// Returns whether the current user is root or not.
func IsRoot() bool {
	if current, err := user.Current(); err == nil {
		if current.Uid == `0` {
			return true
		}
	}

	return false
}

// Returns the first argument if the current user is root, and the second if not.
func RootOr(ifRoot interface{}, notRoot interface{}) interface{} {
	if IsRoot() {
		return ifRoot
	} else {
		return notRoot
	}
}

// The same as RootOr, but returns a string.
func RootOrString(ifRoot interface{}, notRoot interface{}) string {
	if IsRoot() {
		return typeutil.String(ifRoot)
	} else {
		return typeutil.String(notRoot)
	}
}
