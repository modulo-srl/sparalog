package writers

import (
	"os"

	"github.com/modulo-srl/sparalog/logs"
)

type FileWriter struct {
	Writer

	//mu sync.Mutex

	filename string
	file     *os.File
}

// NewFileWriter returns a fileWriter.
func NewFileWriter(filename string) (*FileWriter, error) {
	w := FileWriter{
		filename: filename,
	}

	var err error

	w.file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	w.StartQueue(100, w.onQueueItem)

	return &w, nil
}

func (w *FileWriter) Write(item *logs.Item) {
	w.Enqueue(item)
}

func (w *FileWriter) onQueueItem(item *logs.Item) error {
	//w.mu.Lock()
	//defer w.mu.Unlock()

	s := item.ToString(true, true)

	_, err := w.file.WriteString(s + "\n")
	return err
}

func (w *FileWriter) Stop() {
	w.file.Close()
}
