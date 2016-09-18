package buflog

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type Formatter interface {
	Format(*Entity) ([]byte, error)
}

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	gray    = 37

	DefaultTimestampFormat = time.RFC3339
)

var (
	baseTimestamp time.Time
	isTerminal    bool
	isColored     bool = true
)

func init() {
	baseTimestamp = time.Now()
}

func miniTS() int {
	return int(time.Since(baseTimestamp) / time.Second)
}

type TextFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool
}

func (f *TextFormatter) Format(entity *Entity) ([]byte, error) {

	//isColorTerminal := isTerminal && (runtime.GOOS != "windows")
	//isColored := (f.ForceColors || isColorTerminal) && !f.DisableColors

	b := &bytes.Buffer{}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}

	for _, entry := range entity.cache {
		if isColored {
			f.printColored(b, entity.Id, entry, timestampFormat)
		} else {
			if !f.DisableTimestamp {
				f.appendKeyValue(b, "time", entry.Time.Format(timestampFormat))
			}
			f.appendKeyValue(b, "level", entry.Level.String())
			if entry.Data != "" {
				f.appendKeyValue(b, "msg", entry.Data)
			}
		}
		b.WriteByte('\n')
	}

	return b.Bytes(), nil
}

func (f *TextFormatter) printColored(b *bytes.Buffer, id string, entry *Entry, timestampFormat string) {
	var levelColor int
	switch entry.Level {
	case DebugLevel:
		levelColor = gray
	case WarnLevel:
		levelColor = yellow
	case ErrorLevel, FatalLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	levelText := strings.ToUpper(entry.Level.String())[0:4]

	if !f.FullTimestamp {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%04d]", levelColor, levelText, miniTS())
	} else {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s]", levelColor, levelText, entry.Time.Format(timestampFormat))
	}
	if id != "" {
		fmt.Fprintf(b, " \x1b[%dm[%s]\x1b[0m", levelColor, id)
	}
	fmt.Fprintf(b, " %-44s ", entry.Data)
}

func needsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.') {
			return false
		}
	}
	return true
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {

	b.WriteString(key)
	b.WriteByte('=')

	switch value := value.(type) {
	case string:
		if needsQuoting(value) {
			b.WriteString(value)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	case error:
		errmsg := value.Error()
		if needsQuoting(errmsg) {
			b.WriteString(errmsg)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	default:
		fmt.Fprint(b, value)
	}

	b.WriteByte(' ')
}
