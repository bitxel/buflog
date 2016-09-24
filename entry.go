package buflog

import "time"

// Entry is the basic unit of one log
type Entry struct {
	Time  time.Time
	Level Level
	Data  string
}

// NewEntry return a new Entry obj
func NewEntry(level Level, content string) *Entry {
	en := &Entry{
		Time:  time.Now(),
		Data:  content,
		Level: level,
	}
	return en
}
