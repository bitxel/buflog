package buflog

import (
	"fmt"
	"sync"
	"time"
)

type Entity struct {
	Id       string
	logger   *Logger
	cache    []*Entry
	size     int
	inBuffer bool
	mu       sync.Mutex
}

func (e *Entity) Start() {
	for {
		time.Sleep(e.logger.BufTimeout)
		if e.inBuffer == false {
			return
		}
		e.Flush()
	}
}

func (e *Entity) Empty() {
	e.cache = make([]*Entry, 0)
	e.inBuffer = false
	e.size = 0
}

func (e *Entity) Flush() {
	e.mu.Lock()
	b, _ := e.logger.Formatter.Format(e)
	e.logger.Write(b)
	e.Empty()
	e.mu.Unlock()
}

func (e *Entity) Log(level Level, content string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cache = append(e.cache, NewEntry(level, content))
	e.size += len(content)
	if e.inBuffer == false {
		e.inBuffer = true
		go e.Start()
	}
}

func NewEntity(logger *Logger, id string) *Entity {
	return &Entity{
		Id:       id,
		logger:   logger,
		inBuffer: false,
	}
}

func (e *Entity) Debugf(format string, args ...interface{}) {
	if DebugLevel >= l.Level {
		e.Log(DebugLevel, fmt.Sprintf(format, args))
	}
}

func (e *Entity) Debug(args ...interface{}) {
	if DebugLevel >= l.Level {
		e.Log(DebugLevel, fmt.Sprint(args))
	}
}

func (e *Entity) Debugln(args ...interface{}) {
	e.Debug(args)
}

func (e *Entity) Warnf(format string, args ...interface{}) {
	if WarnLevel >= l.Level {
		e.Log(WarnLevel, fmt.Sprintf(format, args))
	}
}

func (e *Entity) Warn(args ...interface{}) {
	if WarnLevel >= l.Level {
		e.Log(WarnLevel, fmt.Sprint(args))
	}
}

func (e *Entity) Warnln(args ...interface{}) {
	e.Warn(args)
}

func (e *Entity) Infof(format string, args ...interface{}) {
	if InfoLevel >= l.Level {
		e.Log(InfoLevel, fmt.Sprintf(format, args))
	}
}

func (e *Entity) Info(args ...interface{}) {
	if InfoLevel >= l.Level {
		e.Log(InfoLevel, fmt.Sprint(args))
	}
}

func (e *Entity) Infoln(args ...interface{}) {
	e.Info(args)
}

func (e *Entity) Printf(format string, args ...interface{}) {
	e.Infof(format, args)
}

func (e *Entity) Print(args ...interface{}) {
	e.Info(args)
}

func (e *Entity) Println(args ...interface{}) {
	e.Info(args)
}

func (e *Entity) Errorf(format string, args ...interface{}) {
	if ErrorLevel >= l.Level {
		e.Log(ErrorLevel, fmt.Sprintf(format, args))
	}
}

func (e *Entity) Error(args ...interface{}) {
	if ErrorLevel >= l.Level {
		e.Log(ErrorLevel, fmt.Sprint(args))
	}
}

func (e *Entity) Errorln(args ...interface{}) {
	e.Error(args)
}
