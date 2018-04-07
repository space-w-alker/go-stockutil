// Standard logging package, batteries included.
package log

import (
	"fmt"
	"os"

	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/op/go-logging"
)

var backend *logging.LogBackend
var formatted logging.Backend
var leveled logging.LeveledBackend

var defaultLogger *logging.Logger
var ModuleName = ``

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
			fmt.Sprintf(`%%{color}%%{level:.4s}%%{color:reset}[%%{id:04d}] [%%{module:16s}] %%{message}`),
		))

		leveled = logging.AddModuleLevel(formatted)
		logging.SetBackend(leveled)

		defaultLogger = logging.MustGetLogger(ModuleName)
		SetLevel(LogLevel)
	}
}

func Debugging() bool {
	return (LogLevel == DEBUG)
}

func Logger() *logging.Logger {
	initLogging()
	return defaultLogger
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

	switch level {
	case PANIC:
		defaultLogger.Panicf(format, args...)
	case FATAL:
		defaultLogger.Fatalf(format, args...)
	case CRITICAL:
		defaultLogger.Criticalf(format, args...)
	case ERROR:
		defaultLogger.Errorf(format, args...)
	case WARNING:
		defaultLogger.Warningf(format, args...)
	case NOTICE:
		defaultLogger.Noticef(format, args...)
	case INFO:
		defaultLogger.Infof(format, args...)
	default:
		defaultLogger.Debugf(format, args...)
	}
}

func Log(level Level, args ...interface{}) {
	initLogging()

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

func Dump(args ...interface{}) {
	for _, arg := range args {
		Log(DEBUG, typeutil.Dump(arg))
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
