package log

import (
	"fmt"
	"io"
	"regexp"

	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/mgutz/ansi"
)

var rxColorExpr = regexp.MustCompile(`(\$\{(?P<color>[^\}]+)\})`) // ${color}, ${color:mod1:mod2}

func csprintf(colorEnabled bool, format string, args ...interface{}) string {
	out := fmt.Sprintf(format, args...)

	for {
		if match := rxutil.Match(rxColorExpr, out); match != nil {
			colorExpr := match.Group(`color`)
			repl := ``

			// only replace with the actual ANSI escape sequences if we're at a tty
			// or if colors have been explicitly enabled, otherwise just remove the sequences
			if colorEnabled {
				repl = ansi.ColorCode(colorExpr)
			}

			out = match.ReplaceGroup(1, repl)
		} else {
			break
		}
	}

	return out
}

func CSprintf(format string, args ...interface{}) string {
	return csprintf(true, format, args...)
}

func CPrintf(format string, args ...interface{}) (int, error) {
	return fmt.Print(CSprintf(format, args...))
}

func CFPrintf(w io.Writer, format string, args ...interface{}) (int, error) {
	return fmt.Fprint(w, CSprintf(format, args...))
}

func CStripf(format string, args ...interface{}) string {
	return csprintf(false, format, args...)
}
