package writers

import (
	"github.com/getsentry/sentry-go"
	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/writers/templates"
)

type sentryWriter struct {
	queue chan sparalog.Item

	worker *templates.Worker
}

// NewSentryWriter returns a sentryWriter.
func NewSentryWriter() (sparalog.Writer, error) {
	w := sentryWriter{}

	w.worker = templates.NewWorker(&w, 100)

	return &w, nil
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
func (w *sentryWriter) Write(item sparalog.Item) {
	w.worker.Enqueue(item)
}

func (w *sentryWriter) ProcessQueueItem(item sparalog.Item) {
	s := item.String(false, false)
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentryLevels[item.Level])
		sentry.CaptureMessage(s)
	})
}

func (w *sentryWriter) Close() {
	w.worker.Close(3)
}
