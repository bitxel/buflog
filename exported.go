package buflog

import (
	"io"
	"runtime"
	"time"
)

// ID returns a Entity obj with Info(), Error(), Fatal() func as log pkg
func ID(logid string) *Entity {
	return l.ID(logid)
}

// SetBufTimeout is an exported func to set buffer timeout
func SetBufTimeout(d time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.BufTimeout = d
}

// SetBufMaxSize sets logger buffer size
func SetBufMaxSize(size int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.BufMaxSize = size
}

// SetLevel sets min Level to print
func SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Level = level
}

// SetFormatter sets the formatter before print
func SetFormatter(formatter Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Formatter = formatter
}

// #################################################
//  following is the export func in golang log pkg
// #################################################

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
