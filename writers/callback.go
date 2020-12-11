package writers

// Simple writer that forward all writes to a callback.
// Useable as stub writer in unit tests.

import (
	"sync"

	"github.com/modulo-srl/sparalog"
)

type callbackWriter struct {
	mu sync.Mutex

	callback CallbackWriterCallback
}

// CallbackWriterCallback define the writer callback.
type CallbackWriterCallback func(sparalog.Item)

// NewCallbackWriter returns a callbackWriter.
func NewCallbackWriter(callback CallbackWriterCallback) sparalog.Writer {
	w := callbackWriter{
		callback: callback,
	}

	return &w
}

func (w *callbackWriter) Write(item sparalog.Item) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.callback(item)
}

func (w *callbackWriter) Close() {}
