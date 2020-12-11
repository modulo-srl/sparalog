package writers

import (
	"fmt"
	"os"
	"sync"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/writers/templates"
)

type stdoutWriter struct {
	templates.Writer

	mu sync.Mutex
}

// NewStdoutWriter returns a stdoutWriter.
func NewStdoutWriter() sparalog.Writer {
	w := stdoutWriter{}

	return &w
}

func (w *stdoutWriter) Write(item sparalog.Item) sparalog.WriterError {
	w.mu.Lock()
	defer w.mu.Unlock()

	s := item.String(true, true)

	if item.Level <= sparalog.WarnLevel {
		fmt.Fprintln(os.Stderr, s)
		return nil
	}

	fmt.Println(s)
	return nil
}

func (w *stdoutWriter) Close() {}
