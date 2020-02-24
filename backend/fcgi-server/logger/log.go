package logger

import (
	"log"
	"os"
)

const (
	LogLevelNone = iota
	LogLevelFatal
	LogLevelError
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
	LogLevelVerbose
)

var (
	loglevel = LogLevelWarning
	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func LogLevel(l int) {
	loglevel = l
}

func Fatal(v ...interface{}) {
	if loglevel < LogLevelFatal {
		return
	}

	logger.Fatal(v...)
}

func Error(v ...interface{}) {
	if loglevel < LogLevelError {
		return
	}

	logger.Print(v...)
}
func Warning(v ...interface{}) {
	if loglevel < LogLevelWarning {
		return
	}

	logger.Print(v...)
}

func Info(v ...interface{}) {
	if loglevel < LogLevelInfo {
		return
	}

	logger.Print(v...)
}

func Debug(v ...interface{}) {
	if loglevel < LogLevelDebug {
		return
	}

	logger.Print(v...)
}

func Verbose(v ...interface{}) {
	if loglevel < LogLevelVerbose {
		return
	}

	logger.Print(v...)
}
