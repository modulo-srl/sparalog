package writers

import (
	"os"
	"sync"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/writers/templates"
)

type fileWriter struct {
	templates.Writer

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

func (w *fileWriter) Write(item sparalog.Item) sparalog.WriterError {
	w.mu.Lock()
	defer w.mu.Unlock()

	s := item.String(true, true)

	_, err := w.file.WriteString(s + "\n")
	if err != nil {
		return w.ErrorItem(err)
	}

	return nil
}

func (w *fileWriter) Close() {
	w.file.Close()
}
