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

func Fatal(str string, v ...interface{}) {
	if loglevel < LogLevelFatal {
		return
	}

	logger.Fatalf(str, v...)
}

func Error(str string, v ...interface{}) {
	if loglevel < LogLevelError {
		return
	}

	logger.Printf(str, v...)
}
func Warning(str string,v ...interface{}) {
	if loglevel < LogLevelWarning {
		return
	}

	logger.Printf(str, v...)
}
func Info(str string, v ...interface{}) {
	if loglevel < LogLevelInfo {
		return
	}

	logger.Printf(str, v...)
}

func Debug(str string, v ...interface{}) {
	if loglevel < LogLevelDebug {
		return
	}

	logger.Printf(str, v...)
}

func Verbose(str string, v ...interface{}) {
	if loglevel < LogLevelVerbose {
		return
	}

	logger.Printf(str, v...)
}
