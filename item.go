package sparalog

import "time"

// Item is the single item of log.
type Item struct {
	Timestamp time.Time
	Level     Level
	Tag       string

	Line string

	StackTrace string
}

func (i Item) String(timestamp, level bool) string {
	var s string
	var prefixed bool

	if timestamp {
		s = time.Now().UTC().Format("2006-01-02 15:04:05.000") + " "
	}

	if level {
		s += LevelsString[i.Level]
		prefixed = true
	}

	if i.Tag != "" {
		if prefixed {
			s += " "
		}

		s += "[" + i.Tag + "]"
		prefixed = true
	}

	if prefixed {
		s += ": "
	}

	s += i.Line

	if i.StackTrace != "" {
		s += "\n" + i.StackTrace + "\n" // add extra blank line
	}

	return s
}
