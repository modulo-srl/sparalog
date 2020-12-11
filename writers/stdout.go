package writers

import (
	"fmt"
	"os"
	"sync"

	"github.com/modulo-srl/sparalog"
)

type stdoutWriter struct {
	mu sync.Mutex
}

// NewStdoutWriter returns a stdoutWriter.
func NewStdoutWriter() sparalog.Writer {
	w := stdoutWriter{}

	return &w
}

func (w *stdoutWriter) Write(item sparalog.Item) {
	w.mu.Lock()
	defer w.mu.Unlock()

	s := item.String(true, true)

	if item.Level <= sparalog.WarnLevel {
		fmt.Fprintln(os.Stderr, s)
		return
	}

	fmt.Println(s)
}

func (w *stdoutWriter) Close() {}
