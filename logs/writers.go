package logs

// Writers constructors wrapper.

import (
	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/writers"
)

// NewStdoutWriter returns a new stdoutWriter.
func NewStdoutWriter() sparalog.Writer {
	return writers.NewStdoutWriter()
}

// NewFileWriter returns a new fileWriter.
func NewFileWriter(filename string) (sparalog.Writer, error) {
	return writers.NewFileWriter(filename)
}

// NewSentryWriter returns a new sentryWriter.
func NewSentryWriter() sparalog.Writer {
	return writers.NewSentryWriter()
}

// NewTelegramWriter returns a new telegramWriter.
func NewTelegramWriter(botAPIKey string, channelID int) sparalog.Writer {
	return writers.NewTelegramWriter(botAPIKey, channelID)
}

// NewCallbackWriter returns a new callbackWriter.
func NewCallbackWriter(callback writers.CallbackWriterCallback) sparalog.Writer {
	return writers.NewCallbackWriter(callback)
}

// NewCallbackAsyncWriter returns a new callbackAsyncWriter.
func NewCallbackAsyncWriter(callback writers.CallbackWriterCallback) sparalog.Writer {
	return writers.NewCallbackAsyncWriter(callback)
}

// NewTCPWriter return a new tcpWriter.
func NewTCPWriter(address string, port int, debug bool, cb writers.StateChangeCallback) (sparalog.Writer, error) {
	return writers.NewTCPWriter(address, port, debug, cb)
}

// NewSyslogWriter returns a new syslogWriter.
func NewSyslogWriter(tag string) sparalog.Writer {
	return writers.NewSyslogWriter(tag)
}
