package slim

import (
	"log"
	"strings"
)

// slim run mode
const (
	debugCode = iota
	releaseCode
	testCode
)

const (
	// DebugMode indicates slim mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates slim mode is release.
	ReleaseMode = "release"
	// TestMode indicates slim mode is test.
	TestMode = "test"
)

var (
	slimMode = debugCode
	modeName = DebugMode
)

// IsDebugging returns true if the framework is running in debug mode.
// Use SetMode(slime.ReleaseMode) to disable debug mode.
func IsDebugging() bool {
	return slimMode == debugCode
}

// SetMode sets slim mode according to input string.
func SetMode(value string) {
	if value == "" {
		value = DebugMode
	}

	switch value {
	case DebugMode:
		slimMode = debugCode
	case ReleaseMode:
		slimMode = releaseCode
	case TestMode:
		slimMode = testCode
	default:
		panic("slim mode unknown: " + value + " (available mode: debug release test)")
	}

	modeName = value
}

// Mode returns currently slim mode.
func Mode() string {
	return modeName
}

func debugPrintf(format string, values ...interface{}) {
	if IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}

		log.Printf("[slim-debug] "+format, values...)
	}
}
