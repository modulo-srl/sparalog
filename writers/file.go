package writers

import (
	"os"
	"sync"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/item"
	"github.com/modulo-srl/sparalog/writers/base"
)

type fileWriter struct {
	base.Writer

	mu sync.Mutex

	filename string
	file     *os.File
}

// NewFileWriter returns a fileWriter.
func NewFileWriter(filename string) (sparalog.Writer, error) {
	w := fileWriter{
		filename: filename,
	}

	var err error

	w.file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (w *fileWriter) Write(i sparalog.Item) {
	w.mu.Lock()
	defer w.mu.Unlock()

	s := i.ToString(true, true)

	_, err := w.file.WriteString(s + "\n")
	if err != nil {
		w.FeedbackItem(item.NewError(0, err))
	}
}

func (w *fileWriter) Close() {
	w.file.Close()
}
