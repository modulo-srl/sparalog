//go:build !windows
// +build !windows

package writers

import (
	"log/syslog"

	"github.com/modulo-srl/sparalog/logs"
)

type SyslogWriter struct {
	Writer

	sys *syslog.Writer
}

// If tag is empty, the os.Args[0] is used.
func NewSyslogWriter(tag string) *SyslogWriter {
	w := SyslogWriter{}

	var err error
	w.sys, err = syslog.New(syslog.LOG_INFO, tag)

	if err != nil {
		panic(err)
	}

	return &w
}

func (w *SyslogWriter) Write(item *logs.Item) {
	s := item.ToString(false, true)

	switch item.Level {
	case logs.DebugLevel:
		w.sys.Debug(s)
	case logs.InfoLevel:
		w.sys.Info(s)
	case logs.WarningLevel:
		w.sys.Warning(s)
	case logs.ErrorLevel:
		w.sys.Err(s)
	case logs.FatalLevel:
		w.sys.Crit(s)
	}
}

func (w *SyslogWriter) Stop() {
	w.sys.Close()
}
