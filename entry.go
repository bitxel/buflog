package buflog

import "time"

type Entry struct {
	Time  time.Time
	Level Level
	Data  string
}

func NewEntry(level Level, content string) *Entry {
	en := &Entry{
		Time:  time.Now(),
		Data:  content,
		Level: level,
	}
	return en
}
