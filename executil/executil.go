package executil

import (
	"github.com/ghetzel/go-stockutil/fileutil"
	"os"
	"path/filepath"
	"strings"
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
