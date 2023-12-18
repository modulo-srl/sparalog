package logs

// Item del logger.

import (
	"time"

	"github.com/modulo-srl/sparalog/env"
)

type Item struct {
	Ts        time.Time
	Timestamp string // renderizzata

	Level Level

	Prefix, Message string

	StackTrace string

	Payload map[string]any
}

// Genera un nuovo item con timestamp corrente ed eventuale stacktrace.
func newItem(level Level, prefix, msg string, stackCallsToSkip int) *Item {
	ts := time.Now()
	item := &Item{
		Ts:        ts,
		Timestamp: ts.UTC().Format("2006-01-02 15:04:05.000"),
		Level:     level,
		Prefix:    prefix,
		Message:   msg,
	}

	if levelsStackTrace[level] {
		item.GenerateStackTrace(1 + stackCallsToSkip)
	}

	return item
}

// GenerateStackTrace assign the stack trace of current position to the item.
func (i *Item) GenerateStackTrace(callsToSkip int) {
	i.StackTrace = env.StackTrace(callsToSkip + 1)
}

// Setta un valore del payload.
func (i *Item) SetPayload(key string, value any) {
	if i.Payload == nil {
		i.Payload = make(map[string]any)
	}

	i.Payload[key] = value
}

// Formatta l'item sottoforma di stringa.
func (i Item) ToString(timestamp, stacktrace bool) string {
	s := ""

	if timestamp {
		s = i.Timestamp + " "
	}

	s += LevelsString[i.Level]

	if i.Prefix != "" {
		s += " [" + i.Prefix + "]"
	}

	s += ": " + i.Message

	if stacktrace && i.StackTrace != "" {
		s += "\n" + i.StackTrace + "\n" // add extra blank line
	}

	return s
}
