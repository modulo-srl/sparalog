package templates

import (
	"sync"
	"time"

	"github.com/modulo-srl/sparalog"
)

// Worker implements a queue worker with fixed buffer size.
type Worker struct {
	doer WorkerDoer

	queue   chan sparalog.Item
	queueWG sync.WaitGroup
}

// WorkerDoer is the interface that the writer should implements.
type WorkerDoer interface {
	ProcessQueueItem(sparalog.Item)
}

// NewWorker returns a new worker for the writer.
func NewWorker(wd WorkerDoer, queueSize int) *Worker {
	w := Worker{
		doer: wd,
	}

	w.queue = make(chan sparalog.Item, queueSize)

	w.queueWG.Add(1)
	go func() {
		for item := range w.queue {
			w.doer.ProcessQueueItem(item)
		}

		w.queueWG.Done()
	}()

	return &w
}

// Enqueue an item and returns immediately,
// or blocks while the internal queue is full.
func (w *Worker) Enqueue(item sparalog.Item) {
	w.queue <- item
}

// Close the queue and wait for in-progress items.
func (w *Worker) Close(timeoutSecs int) {
	close(w.queue)
	waitTimeout(&w.queueWG, time.Second*time.Duration(timeoutSecs))
}

// Wait for a WaitGroup with a timeout.
// Returns false when timeouted.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		return false
	}
}
