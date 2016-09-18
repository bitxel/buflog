package buflog

import (
	"io"
	"os"
	"sync"
	"time"
)

type Logger struct {
	Out        io.Writer
	Formatter  Formatter
	BufTimeout time.Duration
	BufMaxSize int
	Level      Level
	entry      map[string]*Entity
	mu         sync.Mutex
}

func (l *Logger) Write(b []byte) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Out.Write(b)
}

func (l *Logger) Id(logid string) *Entity {
	if entry, ok := l.entry[logid]; ok {
		return entry
	}
	entry := NewEntity(l, logid)
	l.entry[logid] = entry
	return entry
}

func NewLogger() *Logger {
	return &Logger{
		Out:        os.Stdout,
		Formatter:  &TextFormatter{},
		Level:      InfoLevel,
		BufTimeout: time.Second,
		BufMaxSize: 9999,
		entry:      make(map[string]*Entity),
	}
}
