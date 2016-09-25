package buflog

import (
	"io"
	"sync"
	"time"
)

// Logger is the core struct, it store the control info,
// such as out channel, formatter, time to buffer, buffer size, etc
type Logger struct {
	Out        io.Writer
	Formatter  Formatter
	BufTimeout time.Duration
	BufMaxSize int
	Level      Level
	entry      map[string]*Entity
	mu         sync.Mutex
	Flag       int
	Prefix     string
}

// Write bytes to writter
func (l *Logger) Write(b []byte) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Out.Write(b)
}

// ID returns entity with specific ID
func (l *Logger) ID(logid string) *Entity {
	if entry, ok := l.entry[logid]; ok {
		return entry
	}
	entry := NewEntity(l, logid)
	l.entry[logid] = entry
	return entry
}

// GetEntityByLogAPI return entity with ID="" used by golang log pkg
func (l *Logger) GetEntityByLogAPI() *Entity {
	entry := NewEntity(l, "")
	entry.useLogAPI = true
	return entry
}

// NewLogger return logger obj
func NewLogger(out io.Writer, prefix string) *Logger {
	return &Logger{
		Out:        out,
		Formatter:  &TextFormatter{TimestampFormat: DefaultTimestampFormat},
		Level:      InfoLevel,
		BufTimeout: time.Second,
		BufMaxSize: 9999,
		Prefix:     prefix,
		entry:      make(map[string]*Entity),
	}
}

// SetFormatter sets formattter
func (l *Logger) SetFormatter(f Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Formatter = f
}

// SetFlag sets logger flag
func (l *Logger) SetFlag(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Flag = flag
	l.Formatter.UpdateTimestampFormat(flag)
}

// SetOutput sets out writter
func (l *Logger) SetOutput(out io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Out = out
}

// SetPrefix sets prefix of log
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Prefix = prefix
}
