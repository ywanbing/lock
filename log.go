package lock

import "log"

// Logger as long as this interface is implemented, the default behavior of log can be replaced
type Logger interface {
	PrintLockUsageTime(format string, v ...interface{})
}

// Log is lock internal default log
type Log struct {
	*log.Logger
}

func (l *Log) PrintLockUsageTime(format string, args ...interface{}) {
	l.Printf(format, args...)
}

func NewLog() Logger {
	return &Log{
		Logger: log.Default(),
	}
}
