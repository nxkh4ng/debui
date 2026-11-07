package logger

import (
	"fmt"
	"os"
)

const (
	LevelInfo    = "\x1b[34m[INFO]\x1b[0m"
	LevelWarn    = "\x1b[33m[WARN]\x1b[0m"
	LevelError   = "\x1b[31m[ERROR]\x1b[0m"
	LevelRequest = "\x1b[35m[REQUEST]\x1b[0m"
)

func Infof(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stdout, fmt.Sprintf("%s %s", LevelInfo, format), args...)
}

func Warnf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf("%s %s", LevelWarn, format), args...)
}

func Errorf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf("%s %s", LevelError, format), args...)
}

func Requestf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stdout, fmt.Sprintf("%s %s", LevelRequest, format), args...)
}
