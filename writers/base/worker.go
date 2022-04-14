package base

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
	FeedbackItem(sparalog.Item)
}

// NewWorker returns a new worker for the writer.
func NewWorker(wd WorkerDoer, queueSize int) *Worker {
	w := Worker{
		doer: wd,
	}

	w.queue = make(chan sparalog.Item, queueSize)

	return &w
}

// Start the worker.
func (w *Worker) Start() {
	w.queueWG.Add(1)
	go func() {
		for item := range w.queue {
			w.doer.ProcessQueueItem(item)
		}

		w.queueWG.Done()
	}()
}

// Enqueue an item and returns immediately,
// or blocks while the internal queue is full.
func (w *Worker) Enqueue(item sparalog.Item) {
	w.queue <- item
}

// Stop the queue and wait for in-progress items.
func (w *Worker) Stop(timeoutSecs int) {
	close(w.queue)

	ch := make(chan struct{})

	go func() {
		w.queueWG.Wait()
		close(ch)
	}()

	select {
	case <-ch:
	case <-time.After(time.Second * time.Duration(timeoutSecs)):
	}
}
