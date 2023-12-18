package writers

import (
	"fmt"
	"os"
	"sync"

	"github.com/modulo-srl/sparalog/logs"
)

type StdoutWriter struct {
	Writer

	mu sync.Mutex
}

// NewStdoutWriter returns a stdoutWriter.
func NewStdoutWriter() *StdoutWriter {
	w := StdoutWriter{}

	return &w
}

func (w *StdoutWriter) Write(item *logs.Item) {
	w.mu.Lock()
	defer w.mu.Unlock()

	s := item.ToString(true, true)

	if item.Level <= logs.WarningLevel {
		fmt.Fprintln(os.Stderr, s)
		return
	}

	fmt.Println(s)
}
