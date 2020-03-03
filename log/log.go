package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	multierror "github.com/hashicorp/go-multierror"
	isatty "github.com/mattn/go-isatty"
	"github.com/op/go-logging"
)

var EnableColorExpressions = func() bool {
	if forceColor := os.Getenv(`FORCE_COLOR`); forceColor != `` {
		return typeutil.Bool(forceColor)
	} else {
		return isatty.IsTerminal(os.Stdout.Fd())
	}
}()

var DefaultInterceptStackDepth int = 5
var SynchronousIntercepts = false

var backend *logging.LogBackend
var formatted logging.Backend
var leveled logging.LeveledBackend
var intercepts sync.Map

var defaultLogger *logging.Logger
var ModuleName = ``

type LogFunc func(args ...interface{})
type FormattedLogFunc func(format string, args ...interface{})
type LogInterceptFunc func(level Level, line string, stack StackItems)

// The LOGLEVEL environment variable has final say over the effective log level
// for all users of this package.
var LogLevel Level = func() Level {
	if v := os.Getenv(`LOGLEVEL`); v != `` {
		return GetLevel(v)
	} else {
		return INFO
	}
}()

func initLogging() {
	if defaultLogger == nil {
		backend = logging.NewLogBackend(os.Stderr, ``, 0)

		formatted = logging.NewBackendFormatter(backend, logging.MustStringFormatter(
			fmt.Sprintf(
				`[%%{time:15:04:05} %%{id:04d}] %%{color:bold}%%{level:.4s}%%{color:reset} %%{message}`,
			),
		))

		leveled = logging.AddModuleLevel(formatted)
		logging.SetBackend(leveled)

		defaultLogger = logging.MustGetLogger(ModuleName)
		SetLevel(LogLevel)
	}
}

// Append a function to be called (asynchronously in its own goroutine, or
// synchronously if SynchronousIntercepts is true) for every line logged.
// Returns a UUID that can be later used to deregister the intercept function.
func AddLogIntercept(fn LogInterceptFunc) string {
	id := stringutil.UUID().String()
	intercepts.Store(id, fn)
	return id
}

// Remove the previously-added log intercept function.
func RemoveLogIntercept(id string) {
	intercepts.Delete(id)
}

func Debugging() bool {
	return (LogLevel == DEBUG)
}

func VeryDebugging(features ...string) bool {
	if Debugging() {
		envFeatures := strings.Split(os.Getenv(`DEBUG`), `,`)

		for _, feature := range features {
			for _, ef := range envFeatures {
				if typeutil.Bool(ef) {
					return true
				} else if strings.ToLower(feature) == strings.ToLower(ef) {
					return true
				}
			}
		}
	}

	return false
}

func Logger() *logging.Logger {
	initLogging()
	return defaultLogger
}

// Set the destination Writer where logs will henceforth be written.
func SetOutput(w io.Writer) {
	initLogging()
	backend.Logger.SetOutput(w)
}

func SetLevelString(level string, modules ...string) {
	SetLevel(GetLevel(level), modules...)
}

func SetLevel(level Level, modules ...string) {
	initLogging()

	if lvl, err := logging.LogLevel(level.String()); err == nil {
		if len(modules) == 0 {
			leveled.SetLevel(lvl, ``)
		} else {
			for _, module := range modules {
				leveled.SetLevel(lvl, module)
			}
		}
	} else {
		fmt.Printf("[INVALID LEVEL %v] ", level)
	}
}

func Logf(level Level, format string, args ...interface{}) {
	initLogging()
	callIntercepts(level, fmt.Sprintf(format, args...), StackTrace(DefaultInterceptStackDepth))

	// only replace with the actual ANSI escape sequences if we're at a tty
	// or if colors have been explicitly enabled, otherwise just remove the sequences
	if EnableColorExpressions {
		log(level, CSprintf(format, args...))
	} else {
		log(level, CStripf(format, args...))
	}
}

func Log(level Level, args ...interface{}) {
	// handle this special case where we are handling a fatal-level unformatted nil value,
	// in which case we don't actually want to end the program.
	//
	// NOTE: this obviates the need for FatalfIf, making Fatal() behave in the same way
	//
	if level == FATAL && len(args) == 1 && args[0] == nil {
		return
	}

	initLogging()
	callIntercepts(level, strings.Join(sliceutil.Stringify(args), ` `), StackTrace(DefaultInterceptStackDepth))
	log(level, args...)
}

func log(level Level, args ...interface{}) {
	switch level {
	case PANIC:
		defaultLogger.Panic(args...)
	case FATAL:
		defaultLogger.Fatal(args...)
	case CRITICAL:
		defaultLogger.Critical(args...)
	case ERROR:
		defaultLogger.Error(args...)
	case WARNING:
		defaultLogger.Warning(args...)
	case NOTICE:
		defaultLogger.Notice(args...)
	case INFO:
		defaultLogger.Info(args...)
	default:
		defaultLogger.Debug(args...)
	}
}

func Critical(args ...interface{}) {
	Log(CRITICAL, args...)
}

func Criticalf(format string, args ...interface{}) {
	Logf(CRITICAL, format, args...)
}

func Debug(args ...interface{}) {
	Log(DEBUG, args...)
}

func Debugf(format string, args ...interface{}) {
	Logf(DEBUG, format, args...)
}

// Pretty-print the given arguments to the log at debug-level.
func Dump(args ...interface{}) {
	for _, arg := range args {
		Log(DEBUG, typeutil.Dump(arg))
	}
}

// Same as Dump, but accepts a format string.
func Dumpf(format string, args ...interface{}) {
	for _, arg := range args {
		Logf(DEBUG, format, typeutil.Dump(arg))
	}
}

// Marshal the arguments as indented JSON and log them at debug-level.
func DumpJSON(args ...interface{}) {
	for _, arg := range args {
		if data, err := json.MarshalIndent(arg, ``, `  `); err == nil {
			Log(DEBUG, string(data))
		} else {
			Logf(DEBUG, "DumpJSON: %v", err)
		}
	}
}

func Error(args ...interface{}) {
	Log(ERROR, args...)
}

func Errorf(format string, args ...interface{}) {
	Logf(ERROR, format, args...)
}

func Fatal(args ...interface{}) {
	Log(FATAL, args...)
}

func Fatalf(format string, args ...interface{}) {
	Logf(FATAL, format, args...)
}

func Info(args ...interface{}) {
	Log(INFO, args...)
}

func Infof(format string, args ...interface{}) {
	Logf(INFO, format, args...)
}

func Notice(args ...interface{}) {
	Log(NOTICE, args...)
}

func Noticef(format string, args ...interface{}) {
	Logf(NOTICE, format, args...)
}

func Panic(args ...interface{}) {
	Log(PANIC, args...)
}

func Panicf(format string, args ...interface{}) {
	Logf(PANIC, format, args...)
}

func Warning(args ...interface{}) {
	Log(WARNING, args...)
}

func Warningf(format string, args ...interface{}) {
	Logf(WARNING, format, args...)
}

func Confirm(prompt string) bool {
	return Confirmf(prompt)
}

// Present a confirmation prompt. The function returns true if the user interactively responds
// with "yes" or "y". Otherwise the function returns false.
func Confirmf(format string, args ...interface{}) bool {
	var response string

	fmt.Printf(format, args...)

	if _, err := fmt.Scanln(&response); err != nil {
		panic(err.Error())
	}

	for _, okay := range []string{
		`y`,
		`yes`,
	} {
		if strings.ToLower(okay) == strings.ToLower(response) {
			return true
		}
	}

	return false
}

// Appends one error to another, allowing for operations that return multiple errors
// to remain compatible within a single-valued context.
func AppendError(base error, err error) error {
	if err == nil {
		return base
	} else {
		return multierror.Append(base, err)
	}
}

// Invoke Fatal() if the given error is not nil.
func FatalIf(err error) {
	if err != nil {
		Fatal(err)
	}
}

// Invoke Fatalf() if the given error is not nil.
func FatalfIf(format string, err error) {
	if err != nil {
		Fatalf(format, err)
	}
}

// call all registered intercept functions using the given arguments.
func callIntercepts(level Level, line string, stack StackItems) {
	intercepts.Range(func(_ interface{}, value interface{}) bool {
		if fn, ok := value.(LogInterceptFunc); ok {
			// for levels CRITICAL and worse, call intercepts synchronously in case we're
			// panicking and about to tear crap down.  Since these intercepts should run BEFORE
			// the log line is emitted, this should ensure the intercept definitely runs before
			// any of that goes down.
			if level <= CRITICAL || SynchronousIntercepts {
				fn(level, line, stack)
			} else {
				go fn(level, line, stack)
			}
		}

		return true
	})
}
