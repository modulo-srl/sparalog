package base

import (
	"crypto/rand"
	"fmt"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/item"
)

// Writer implements the base methods.
type Writer struct {
	id sparalog.WriterID

	feedbackCh chan sparalog.Item
}

func (w *Writer) Open() error { return nil }
func (w *Writer) Close()      {}

func (w *Writer) ID() sparalog.WriterID {
	if w.id == "" {
		bb := make([]byte, 8)
		rand.Read(bb)
		w.id = sparalog.WriterID(fmt.Sprintf("%X", bb))
	}

	return w.id
}

// SetFeedbackChan set a channel to the level default writer of the logger.
func (w *Writer) SetFeedbackChan(ch chan sparalog.Item) {
	w.feedbackCh = ch
}

// Feedback generate an item and send it to the level default writer of the logger.
func (w *Writer) Feedback(level sparalog.Level, args ...interface{}) {
	i := item.New(level, fmt.Sprint(args...))

	w.FeedbackItem(i)
}

// Feedbackf generate an item and send it to the level default writer of the logger.
func (w *Writer) Feedbackf(level sparalog.Level, format string, args ...interface{}) {
	i := item.New(level, fmt.Sprintf(format, args...))

	w.FeedbackItem(i)
}

// FeedbackItem send an item to the level default writer of ther logger.
func (w *Writer) FeedbackItem(item sparalog.Item) {
	if w.feedbackCh == nil {
		return
	}

	w.feedbackCh <- item
}
