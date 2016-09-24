package buflog

import (
	"io"
	"os"
	"runtime"
)

const (
	_ = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

// These flags are copied in order to be compatible with golang log pkg
const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

var (
	l = NewLogger(os.Stdout, "")
)

// ID returns a Entity obj with Info(), Error(), Fatal() func as log pkg
func ID(logid string) *Entity {
	return l.ID(logid)
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
	case PanicLevel:
		return "panic"
	}

	return "unknown"
}

// Fatal prints the log then exits
func Fatal(v ...interface{}) {
	l.GetEntityByLogAPI().Fatal(v...)
}

// Fatalf prints formated log then exits
func Fatalf(format string, v ...interface{}) {
	l.GetEntityByLogAPI().Fatalf(format, v...)
}

// Fatalln works as the Fatal()
func Fatalln(v ...interface{}) {
	l.GetEntityByLogAPI().Fatalln(v...)
}

// Flags returns the flag defined
func Flags() int {
	return l.Flag
}

// Output prints the caller
func Output(calldepth int, s string) error {
	if l.Flag&(Lshortfile|Llongfile) != 0 {
		_, file, line, ok := runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		if l.Flag&Lshortfile != 0 {
			file = fileNameLongToShort(file)
		}
		l.GetEntityByLogAPI().Printf("%s:%d %s", file, line, s)
	}
	return nil
}

// Panic prints the log and panic
func Panic(v ...interface{}) {
	l.GetEntityByLogAPI().Panic(v...)
}

// Panicf prints formated log and panic
func Panicf(format string, v ...interface{}) {
	l.GetEntityByLogAPI().Panicf(format, v...)
}

// Panicln is the same as panic() in buflog
func Panicln(v ...interface{}) {
	l.GetEntityByLogAPI().Panicln(v...)
}

// Prefix returns log prefix
func Prefix() string {
	return l.Prefix
}

// Print prints log
func Print(v ...interface{}) {
	l.GetEntityByLogAPI().Print(v...)
}

// Printf prints formated log
func Printf(format string, v ...interface{}) {
	l.GetEntityByLogAPI().Printf(format, v...)
}

// Println is the same as print()
func Println(v ...interface{}) {
	l.GetEntityByLogAPI().Println(v...)
}

// SetFlags set options when print log
// check the constant part of log pkg
// https://golang.org/pkg/log/
func SetFlags(flag int) {
	l.SetFlag(flag)
}

// SetOutput defines the output writer, default os.StdErr
func SetOutput(w io.Writer) {
	l.SetOutput(w)
}

// SetPrefix sets the content at the beginning of log
func SetPrefix(prefix string) {
	l.SetPrefix(prefix)
}

func fileNameLongToShort(file string) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			return short
		}
	}
	return file
}
