package logs

// Items constructors wrapper.

import (
	"fmt"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logger"
)

// NewItem generate a new log item for the default logger.
func NewItem(level sparalog.Level, args ...string) sparalog.Item {
	return Default.NewItem(level, args...)
}

// NewItemf generate a new log item for the default logger.
func NewItemf(level sparalog.Level, format string, args ...string) sparalog.Item {
	return Default.NewItemf(level, format, args...)
}

// NewError generate a new log item wrapping an error.
func NewError(err error) sparalog.Item {
	line := err.Error() // as logger.NewError()
	return Default.(*logger.Logger).NewErrorItem(line)
}

// NewErrorf generate a new log item wrapping an error, as Errorf().
func NewErrorf(format string, a ...interface{}) sparalog.Item {
	line := fmt.Errorf(format, a...).Error() // as logger.NewErrorf()
	return Default.(*logger.Logger).NewErrorItem(line)
}
