package writers

import (
	"github.com/getsentry/sentry-go"
	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/writers/templates"
)

type sentryWriter struct {
	templates.Writer

	queue chan sparalog.Item

	worker *templates.Worker
}

// NewSentryWriter returns a sentryWriter.
func NewSentryWriter() sparalog.Writer {
	w := sentryWriter{}

	w.worker = templates.NewWorker(&w, 100)

	return &w
}

var sentryLevels = [sparalog.LevelsCount]sentry.Level{
	sentry.LevelFatal,
	sentry.LevelError,
	sentry.LevelWarning,
	sentry.LevelInfo,
	sentry.LevelDebug,
	sentry.LevelDebug,
}

// Write enqueue an item and returns immediately,
// or blocks while the internal queue is full.
func (w *sentryWriter) Write(item sparalog.Item) sparalog.WriterError {
	w.worker.Enqueue(item)
	return nil
}

func (w *sentryWriter) ProcessQueueItem(item sparalog.Item) sparalog.WriterError {
	s := item.String(false, false)

	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentryLevels[item.Level])
		sentry.CaptureMessage(s)
	})

	return nil
}

func (w *sentryWriter) Close() {
	w.worker.Close(3)
}
