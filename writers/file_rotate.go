package writers

import (
	"compress/gzip"
	"io"
	"os"
	"time"

	"github.com/modulo-srl/sparalog/logs"
)

type FileRotateWriter struct {
	Writer

	lastRotation time.Time
	rotateAfter  time.Duration

	deleteNotCritical bool
	critical          bool

	filename string
	file     *os.File
}

// NewFileRotateWriter returns a FileRotateWriter.
// - deleteNotCritical: if True removes the rotations that do not contain critical logs.
func NewFileRotateWriter(filename string, rotateAfter time.Duration, deleteNotCritical bool) (*FileRotateWriter, error) {
	w := FileRotateWriter{
		filename:          filename,
		rotateAfter:       rotateAfter,
		lastRotation:      time.Now(),
		deleteNotCritical: deleteNotCritical,
	}

	var err error

	w.file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	w.StartQueue(100, w.onQueueItem)

	return &w, nil
}

func (w *FileRotateWriter) Write(item *logs.Item) {
	w.Enqueue(item)
}

func (w *FileRotateWriter) onQueueItem(item *logs.Item) error {
	s := item.ToString(true, true)

	_, err := w.file.WriteString(s + "\n")
	if err != nil {
		return err
	}

	if item.Level <= logs.WarningLevel {
		w.critical = true
	}

	now := time.Now()
	if now.Sub(w.lastRotation) > w.rotateAfter {
		err = w.rotate()
		if err != nil {
			return err
		}

		w.lastRotation = now
		w.critical = false
	}

	return nil
}

func (w *FileRotateWriter) Stop() {
	w.file.Close()
}

func (w *FileRotateWriter) rotate() error {
	var err error

	w.file.Close()

	if !w.deleteNotCritical || w.critical {
		err = w.compress()
		if err != nil {
			return err
		}
	}

	w.file, err = os.OpenFile(w.filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (w *FileRotateWriter) compress() error {
	fr, err := os.Open(w.filename)
	if err != nil {
		return err
	}
	defer fr.Close()

	fi, err := fr.Stat()
	if err != nil {
		return err
	}
	if fi.Size() == 0 {
		// If the file is empty, avoid compressing it.
		return nil
	}

	newFn := w.filename + "." + time.Now().Format("2006-01-02_03-04-05") + ".gz"
	fw, err := os.Create(newFn)
	if err != nil {
		return err
	}
	defer fw.Close()

	zip := gzip.NewWriter(fw)
	defer zip.Close()

	_, err = io.Copy(zip, fr)
	if err != nil {
		return err
	}

	err = zip.Flush()
	if err != nil {
		return err
	}

	return nil
}
