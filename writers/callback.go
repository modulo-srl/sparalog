package writers

// Simple writer that forward all writes to a callback.
// Useable as stub writer in unit tests.

import (
	"sync"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/writers/templates"
)

type callbackWriter struct {
	templates.Writer

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

func (w *callbackWriter) Write(item sparalog.Item) sparalog.WriterError {
	w.mu.Lock()
	defer w.mu.Unlock()

	err := w.callback(item)
	if err != nil {
		return w.ErrorItem(err)
	}

	return nil
}

func (w *callbackWriter) Close() {}

type callbackAsyncWriter struct {
	templates.Writer

	worker *templates.Worker

	callback CallbackWriterCallback
}

// NewCallbackAsyncWriter returns a callbackAsyncWriter.
func NewCallbackAsyncWriter(callback CallbackWriterCallback) sparalog.Writer {
	w := callbackAsyncWriter{
		callback: callback,
	}

	w.worker = templates.NewWorker(&w, 100)

	return &w
}

func (w *callbackAsyncWriter) Close() {
	w.worker.Close(1)
}

func (w *callbackAsyncWriter) Write(item sparalog.Item) sparalog.WriterError {
	w.worker.Enqueue(item)
	return nil
}

func (w *callbackAsyncWriter) ProcessQueueItem(item sparalog.Item) sparalog.WriterError {
	err := w.callback(item)
	if err != nil {
		return w.ErrorItem(err)
	}

	return nil
}
