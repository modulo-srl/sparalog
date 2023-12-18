package writers

// Simple writer that forward all writes to a callback.
// Useable as stub writer in unit tests.

import (
	"sync"

	"github.com/modulo-srl/sparalog/logs"
)

type CallbackWriter struct {
	Writer

	mu sync.Mutex

	callback CallbackWriterCallback
}

// CallbackWriterCallback define the writer callback.
type CallbackWriterCallback func(*logs.Item) error

// NewCallbackWriter returns a callbackWriter.
func NewCallbackWriter(callback CallbackWriterCallback) *CallbackWriter {
	w := CallbackWriter{
		callback: callback,
	}

	return &w
}

func (w *CallbackWriter) Write(item *logs.Item) {
	w.mu.Lock()
	defer w.mu.Unlock()

	err := w.callback(item)
	if err != nil {
		w.FeedbackError(err)
		return
	}
}

type CallbackAsyncWriter struct {
	Writer

	callback CallbackWriterCallback
}

// NewCallbackAsyncWriter returns a callbackAsyncWriter.
func NewCallbackAsyncWriter(callback CallbackWriterCallback) *CallbackAsyncWriter {
	return &CallbackAsyncWriter{
		callback: callback,
	}
}

func (w *CallbackAsyncWriter) Start() error {
	w.StartQueue(100, w.onQueueItem)
	return nil
}

func (w *CallbackAsyncWriter) Stop() {
	w.StopQueue(1)
}

func (w *CallbackAsyncWriter) Write(item *logs.Item) {
	w.Enqueue(item)
}

func (w *CallbackAsyncWriter) onQueueItem(item *logs.Item) error {
	return w.callback(item)
}
