package writers

import (
	"log/syslog"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/writers/base"
)

type syslogWriter struct {
	base.Writer

	sys *syslog.Writer
}

// If tag is empty, the os.Args[0] is used.
func NewSyslogWriter(tag string) sparalog.Writer {
	w := syslogWriter{}

	var err error
	w.sys, err = syslog.New(syslog.LOG_INFO, tag)

	if err != nil {
		panic(err)
	}

	return &w
}

func (w *syslogWriter) Write(item sparalog.Item) {
	switch item.Level() {
	case sparalog.TraceLevel:
		w.sys.Debug(item.ToString(false, false))
	case sparalog.DebugLevel:
		w.sys.Debug(item.ToString(false, false))
	case sparalog.InfoLevel:
		w.sys.Info(item.ToString(false, false))
	case sparalog.WarnLevel:
		w.sys.Warning(item.ToString(false, false))
	case sparalog.ErrorLevel:
		w.sys.Err(item.ToString(false, false))
	case sparalog.FatalLevel:
		w.sys.Crit(item.ToString(false, false))
	}
}

func (w *syslogWriter) Close() {
	w.sys.Close()
}
