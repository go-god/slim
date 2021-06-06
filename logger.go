package slim

import (
	"strings"
)

// Logger is logger interface.
type Logger interface {
	Printf(format string, args ...interface{})
}

// LoggerFunc is a bridge between Logger and any third party logger.
type LoggerFunc func(format string, args ...interface{})

// Printf implements Logger interface.
func (f LoggerFunc) Printf(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}

	f(format, args...)
}

// dummy logger writes nothing.
var dummyLogger = LoggerFunc(func(string, ...interface{}) {})

var (
	// LogEntry default logger entry
	LogEntry Logger = dummyLogger
)
