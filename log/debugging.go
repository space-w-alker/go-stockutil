package log

import (
	"fmt"
	"runtime"
	"strings"
)

type StackItem struct {
	ProgramCounter uintptr
	Filename       string
	Line           int
	Function       string
}

var MaxStackTraceDepth = 32

func (self StackItem) String() string {
	var line []string

	if self.Function != `` {
		line = append(line, fmt.Sprintf("function %v", self.Function))
	}

	if self.Filename != `` && self.Line > 0 {
		line = append(line, fmt.Sprintf("%s, line %d", self.Filename, self.Line))
	}

	return strings.Join(line, "\n")
}

// Retrieves details about the call stack that led to this function call.
func StackTrace(skip int) []StackItem {
	pc := make([]uintptr, MaxStackTraceDepth)
	items := make([]StackItem, 0)

	if n := runtime.Callers(skip, pc); n <= len(pc) {
		pc = pc[:n]

		if frames := runtime.CallersFrames(pc); frames != nil {
			for i := 0; len(items) <= len(pc); i++ {
				frame, more := frames.Next()

				items = append(items, StackItem{
					ProgramCounter: frame.PC,
					Function:       frame.Function,
					Filename:       frame.File,
					Line:           frame.Line,
				})

				if !more {
					break
				}
			}
		}
	}

	return items
}

// Logs the current stack trace as debug log output.
func DebugStack() {
	Debug("Stack trace:")

	for i, item := range StackTrace(3) {
		for j, line := range strings.Split(item.String(), "\n") {
			if j == 0 {
				Debugf("  % 2d: %v", i, line)
			} else {
				Debugf("          %v", line)
			}
		}
	}
}
