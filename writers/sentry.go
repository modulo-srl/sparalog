package writers

import (
	"github.com/getsentry/sentry-go"
	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/writers/base"
)

type sentryWriter struct {
	base.Writer

	queue chan sparalog.Item

	worker *base.Worker
}

// NewSentryWriter returns a sentryWriter.
func NewSentryWriter() sparalog.Writer {
	w := sentryWriter{}

	w.worker = base.NewWorker(&w, 100)

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
func (w *sentryWriter) Write(item sparalog.Item) {
	w.worker.Enqueue(item)
}

func (w *sentryWriter) ProcessQueueItem(item sparalog.Item) {
	s := item.ToString(false, false)

	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentryLevels[item.Level()])

		scope.SetTags(item.Tags())
		scope.SetContexts(item.Data())

		sentry.CaptureMessage(s)
	})
}

func (w *sentryWriter) Open() error {
	w.worker.Start()
	return nil
}

func (w *sentryWriter) Close() {
	w.worker.Stop(3)
}
