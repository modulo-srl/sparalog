//go:build !windows
// +build !windows

package logs

// Writers constructors wrapper.

import (
	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/writers"
)

// NewSyslogWriter returns a new syslogWriter.
func NewSyslogWriter(tag string) sparalog.Writer {
	return writers.NewSyslogWriter(tag)
}
