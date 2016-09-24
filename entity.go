package buflog

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

// Entity is the group of entries with specific Id
type Entity struct {
	ID        string
	logger    *Logger
	cache     []*Entry
	size      int
	inBuffer  bool
	mu        sync.Mutex
	useLogAPI bool
}

// Start a goroutine, and flush the log when time is out or buffer full
func (e *Entity) Start() {
	for {
		time.Sleep(e.logger.BufTimeout)
		if e.inBuffer == false {
			return
		}
		e.Flush()
	}
}

// Empty the buffer
func (e *Entity) Empty() {
	e.cache = make([]*Entry, 0)
	e.inBuffer = false
	e.size = 0
}

// Flush the log to out channel immediately
func (e *Entity) Flush() {
	e.mu.Lock()
	b, _ := e.logger.Formatter.Format(e)
	e.logger.Write(b)
	e.Empty()
	e.mu.Unlock()
}

// Log is the main function to create an entry
func (e *Entity) Log(level Level, content string) {

	if e.logger.Flag&(Lshortfile|Llongfile) != 0 {
		depth := 2
		if e.useLogAPI {
			depth = 3
		}
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			file = "???"
			line = 0
		}
		if l.Flag&Lshortfile != 0 {
			file = fileNameLongToShort(file)
		}
		content = fmt.Sprintf("%s:%d %s", file, line, content)
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.cache = append(e.cache, NewEntry(level, content))
	e.size += len(content)
	if e.inBuffer == false {
		e.inBuffer = true
		if level >= FatalLevel {
			// Flush all, just wait until timeout, not a perfect way
			time.Sleep(e.logger.BufTimeout)
			e.Flush()
		} else {
			go e.Start()
		}
	}
}

// NewEntity return a new obj
func NewEntity(logger *Logger, id string) *Entity {
	return &Entity{
		ID:       id,
		logger:   logger,
		inBuffer: false,
	}
}

// Debugf returns the formated debug level log
func (e *Entity) Debugf(format string, args ...interface{}) {
	if DebugLevel >= l.Level {
		e.Log(DebugLevel, fmt.Sprintf(format, args...))
	}
}

// Debug prints debug level log
func (e *Entity) Debug(args ...interface{}) {
	if DebugLevel >= l.Level {
		e.Log(DebugLevel, fmt.Sprint(args...))
	}
}

// Debugln works the same as debug()
func (e *Entity) Debugln(args ...interface{}) {
	if DebugLevel >= l.Level {
		e.Log(DebugLevel, fmt.Sprint(args...))
	}
}

// Warnf prints formated warn level log
func (e *Entity) Warnf(format string, args ...interface{}) {
	if WarnLevel >= l.Level {
		e.Log(WarnLevel, fmt.Sprintf(format, args...))
	}
}

// Warn prints warn level log
func (e *Entity) Warn(args ...interface{}) {
	if WarnLevel >= l.Level {
		e.Log(WarnLevel, fmt.Sprint(args...))
	}
}

// Warnln works the same as warn()
func (e *Entity) Warnln(args ...interface{}) {
	e.Warn(args...)
}

func (e *Entity) Infof(format string, args ...interface{}) {
	if InfoLevel >= l.Level {
		e.Log(InfoLevel, fmt.Sprintf(format, args...))
	}
}

func (e *Entity) Info(args ...interface{}) {
	if InfoLevel >= l.Level {
		e.Log(InfoLevel, fmt.Sprint(args...))
	}
}

func (e *Entity) Infoln(args ...interface{}) {
	e.Info(args...)
}

func (e *Entity) Printf(format string, args ...interface{}) {
	if InfoLevel >= l.Level {
		e.Log(InfoLevel, fmt.Sprintf(format, args...))
	}
}

func (e *Entity) Print(args ...interface{}) {
	if InfoLevel >= l.Level {
		e.Log(InfoLevel, fmt.Sprint(args...))
	}
}

func (e *Entity) Println(args ...interface{}) {
	e.Info(args...)
}

func (e *Entity) Errorf(format string, args ...interface{}) {
	if ErrorLevel >= l.Level {
		e.Log(ErrorLevel, fmt.Sprintf(format, args...))
	}
}

func (e *Entity) Error(args ...interface{}) {
	if ErrorLevel >= l.Level {
		e.Log(ErrorLevel, fmt.Sprint(args...))
	}
}

func (e *Entity) Errorln(args ...interface{}) {
	e.Error(args...)
}

func (e *Entity) Fatalf(format string, args ...interface{}) {
	e.Log(FatalLevel, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (e *Entity) Fatal(args ...interface{}) {
	e.Log(FatalLevel, fmt.Sprint(args...))
	os.Exit(1)
}

func (e *Entity) Fatalln(args ...interface{}) {
	e.Fatal(args...)
}

func (e *Entity) Panicf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	e.Log(PanicLevel, s)
	panic(s)
}

func (e *Entity) Panic(args ...interface{}) {
	s := fmt.Sprint(args...)
	e.Log(PanicLevel, s)
	panic(s)
}

func (e *Entity) Panicln(args ...interface{}) {
	e.Panic(args...)
}
