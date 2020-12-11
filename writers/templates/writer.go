package templates

import (
	"time"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/env"
)

type Writer struct {
	feedbackCh chan sparalog.WriterError
}

// SetFeedbackChan set a channel to the default level writer of ther logger.
func (w *Writer) SetFeedbackChan(ch chan sparalog.WriterError) {
	w.feedbackCh = ch
}

// FeedbackError send ad error to the default level writer of ther logger.
func (w *Writer) FeedbackError(err sparalog.WriterError) {
	if w.feedbackCh == nil {
		return
	}

	w.feedbackCh <- err
}

func (w *Writer) ErrorItem(err error) sparalog.WriterError {
	e := sparalog.Item{
		Timestamp:  time.Now(),
		Level:      sparalog.ErrorLevel,
		Line:       err.Error(),
		StackTrace: env.StackTrace(3),
	}

	return &e
}
