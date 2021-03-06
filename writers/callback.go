package writers

// Simple writer that forward all writes to a callback.
// Useable as stub writer in unit tests.

import (
	"sync"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/item"
	"github.com/modulo-srl/sparalog/writers/base"
)

type callbackWriter struct {
	base.Writer

	mu sync.Mutex

	callback CallbackWriterCallback
}

// CallbackWriterCallback define the writer callback.
type CallbackWriterCallback func(sparalog.Item) error

// NewCallbackWriter returns a callbackWriter.
func NewCallbackWriter(callback CallbackWriterCallback) sparalog.Writer {
	w := callbackWriter{
		callback: callback,
	}

	return &w
}

func (w *callbackWriter) Write(i sparalog.Item) {
	w.mu.Lock()
	defer w.mu.Unlock()

	err := w.callback(i)
	if err != nil {
		w.FeedbackItem(item.NewError(1, err))
		return
	}
}

func (w *callbackWriter) Close() {}

type callbackAsyncWriter struct {
	base.Writer

	worker *base.Worker

	callback CallbackWriterCallback
}

// NewCallbackAsyncWriter returns a callbackAsyncWriter.
func NewCallbackAsyncWriter(callback CallbackWriterCallback) sparalog.Writer {
	w := callbackAsyncWriter{
		callback: callback,
	}

	w.worker = base.NewWorker(&w, 100)

	return &w
}

func (w *callbackAsyncWriter) Close() {
	w.worker.Close(1)
}

func (w *callbackAsyncWriter) Write(item sparalog.Item) {
	w.worker.Enqueue(item)
}

func (w *callbackAsyncWriter) ProcessQueueItem(i sparalog.Item) {
	err := w.callback(i)
	if err != nil {
		w.FeedbackItem(item.NewError(1, err))
	}
}
