package buflog

import ()

const (
	_ = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var l = NewLogger()

func Id(logid string) *Entity {
	return l.Id(logid)
}

// Level type
type Level uint8

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
		//case PanicLevel:
		//	return "panic"
	}

	return "unknown"
}
