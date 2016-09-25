package buflog

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// Formatter is the interface to Format an log entry
type Formatter interface {
	Format(*Entity) ([]byte, error)
	UpdateTimestampFormat(int)
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
	level2color   []int
	isColored     = true
)

func init() {
	baseTimestamp = time.Now()
	level2color = make([]int, 7)
	level2color[DebugLevel] = gray
	level2color[InfoLevel] = blue
	level2color[WarnLevel] = yellow
	level2color[ErrorLevel] = red
	level2color[FatalLevel] = red
	level2color[PanicLevel] = red
}

func miniTS() int {
	return int(time.Since(baseTimestamp) / time.Second)
}

// TextFormatter is an implementation of Formatter interface
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
	IsUTC          bool
}

// UpdateTimestampFormat by flag parameter
func (f *TextFormatter) UpdateTimestampFormat(flag int) {
	format := ""
	if flag&LUTC != 0 {
		f.IsUTC = true
	}
	if flag&Ldate != 0 {
		format += "2006/01/02 "
	}
	if l.Flag&(Ltime|Lmicroseconds) != 0 {
		format += "15:04:05"
		if flag&Lmicroseconds != 0 {
			format += ".999999"
		}
	}
	f.TimestampFormat = format
}

// Format return the formated entity data to bytes, return it and error
func (f *TextFormatter) Format(entity *Entity) ([]byte, error) {

	//isColorTerminal := isTerminal && (runtime.GOOS != "windows")
	//isColored := (f.ForceColors || isColorTerminal) && !f.DisableColors

	b := &bytes.Buffer{}

	for _, entry := range entity.cache {

		if entity.logger.Prefix != "" {
			b.WriteString(entity.logger.Prefix)
		}

		if f.IsUTC {
			entry.Time = entry.Time.UTC()
		}

		// TODO may encounter effective issue using so many fmt
		f.appendKeyValue(b, level2color[entry.Level], entry.Level, "level", strings.ToUpper(entry.Level.String())[0:4])
		if !f.DisableTimestamp {
			if f.FullTimestamp {
				f.appendKeyValue(b, nocolor, entry.Level, "time", fmt.Sprintf("[%s]", entry.Time.Format(f.TimestampFormat)))
			} else {
				f.appendKeyValue(b, nocolor, entry.Level, "time", fmt.Sprintf("[%d]", miniTS()))
			}
		}
		if entity.ID != "" {
			f.appendKeyValue(b, green, entry.Level, "id", entity.ID)
		}
		if entity.logger.Flag&(Llongfile|Lshortfile) != 0 {
			f.appendKeyValue(b, nocolor, entry.Level, "file", "")
		}
		f.appendKeyValue(b, nocolor, entry.Level, "msg", entry.Data)
		b.WriteByte('\n')
	}

	return b.Bytes(), nil
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

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, color int, level Level, key string, value interface{}) {
	if isColored {
		if color != 0 {
			fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m", color, value)
		} else {
			fmt.Fprint(b, value)
		}
	} else {
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
	}
	b.WriteByte(' ')
}
