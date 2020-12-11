package logs

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

// NewCallbackWriter returns a new callbackWriter.
func NewCallbackWriter(callback writers.CallbackWriterCallback) sparalog.Writer {
	return writers.NewCallbackWriter(callback)
}

// NewSentryWriter returns a new sentryWriter.
func NewSentryWriter() (sparalog.Writer, error) {
	return writers.NewSentryWriter()
}

// NewTelegramWriter returns a new telegramWriter.
func NewTelegramWriter(botAPIKey string, channelID int) (sparalog.Writer, error) {
	return writers.NewTelegramWriter(botAPIKey, channelID)
}
