package item

import (
	"fmt"
	"time"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/env"
)

// New generates a new item with current timestamp.
func New(level sparalog.Level, line string) sparalog.Item {
	i := Item{
		timestamp: time.Now(),
		level:     level,
		line:      line,
	}

	return &i
}

// NewError generates a new item containing the stack trace.
// If skipStackCalls = -1 no stacktrace will be produced at all.
func NewError(skipStackCalls int, err error) sparalog.Item {
	i := New(sparalog.ErrorLevel, err.Error())

	if skipStackCalls >= 0 {
		i.GenerateStackTrace(skipStackCalls + 1)
	}

	return i
}

// NewErrorf generate a new log item wrapping an error, as Errorf().
// If skipStackCalls = -1 no stacktrace will be produced at all.
func NewErrorf(skipStackCalls int, format string, a ...interface{}) sparalog.Item {
	i := New(sparalog.ErrorLevel, fmt.Errorf(format, a...).Error())

	if skipStackCalls >= 0 {
		i.GenerateStackTrace(skipStackCalls + 1)
	}

	return i
}

// Item implements sparalog.Item
type Item struct {
	Context

	logger sparalog.ItemLogger

	timestamp time.Time

	level sparalog.Level

	line string

	stackTrace string
}

// SetLogger assigns a logger to the item, for Log() further using.
func (i *Item) SetLogger(logger sparalog.ItemLogger) {
	i.logger = logger
}

// Level sets the level of the item.
func (i *Item) Level() sparalog.Level {
	return i.level
}

// Line gets the desc of the item.
func (i *Item) Line() string {
	return i.line
}

// SetLine sets de desc of the item.
func (i *Item) SetLine(s string) {
	i.line = s
}

// StackTrace gets the stack trace of the item.
func (i *Item) StackTrace() string {
	return i.stackTrace
}

// SetStackTrace assign a stack trace to the item.
func (i *Item) SetStackTrace(s string) {
	i.stackTrace = s
}

// GenerateStackTrace assign the stack trace of current position to the item.
func (i *Item) GenerateStackTrace(callsToSkip int) {
	i.SetStackTrace(env.StackTrace(callsToSkip + 1))
}

// Log sends the item to the assigned logger.
func (i *Item) Log() {
	if i.logger == nil {
		return
	}

	i.logger.LogItem(i)
}

// ToString convert the item to string using a standard format.
func (i Item) ToString(timestamp, level bool) string {
	var s string
	var prefixed bool

	if timestamp {
		s = time.Now().UTC().Format("2006-01-02 15:04:05.000") + " "
	}

	if level {
		s += sparalog.LevelsString[i.level]
		prefixed = true
	}

	prefix := i.Context.RenderPrefix()
	if prefix != "" {
		if prefixed {
			s += " "
		}

		s += "[" + prefix + "]"
		prefixed = true
	}

	if prefixed {
		s += ": "
	}

	s += i.line

	if i.stackTrace != "" {
		s += "\n" + i.stackTrace + "\n" // add extra blank line
	}

	return s
}
