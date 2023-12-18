package writers

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/modulo-srl/sparalog/logs"
)

// Implementa i metodi base.
type Writer struct {
	id string

	feedbackCh chan *logs.Item

	queue   chan *logs.Item
	queueWG sync.WaitGroup
}

func (w *Writer) ID() string {
	if w.id == "" {
		bb := make([]byte, 8)
		rand.Read(bb)
		w.id = fmt.Sprintf("%X", bb)
	}

	return w.id
}

func (w *Writer) Start() error { return nil }
func (w *Writer) Stop()        {}

type OnItemFunc func(*logs.Item) error

// Avvia il worker di gestione della coda,
// invocando la callback OnItemFunc() per ogni item da processare.
// Se la callback ritorna errore questo viene feedbackato al writer di default
// del rispettivo livello.
func (w *Writer) StartQueue(queueSize int, f OnItemFunc) {
	w.queue = make(chan *logs.Item, queueSize)

	w.queueWG.Add(1)
	go func() {
		for item := range w.queue {
			err := f(item)
			if err != nil {
				//fmt.Println(err)
				w.FeedbackError(err)
			}
		}

		w.queueWG.Done()
	}()
}

// Ritorna immediatamente.
func (w *Writer) Enqueue(item *logs.Item) {
	w.queue <- item
}

// Finisce di consegnare gli item rimanenti in coda e termina.
func (w *Writer) StopQueue(timeoutSecs int) {
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

// Imposta il canale interno di feeback.
// Viene invocata dal logger quando imposta un nuovo writer per un certo livello.
func (w *Writer) SetFeedbackChan(ch chan *logs.Item) {
	w.feedbackCh = ch
}

// Genera un item e lo invia al writer di default del rispettivo livello.
func (w *Writer) Feedback(level logs.Level, args ...any) {
	if w.feedbackCh == nil {
		return
	}

	w.feedbackCh <- logs.NewItem(level, w.ID(), fmt.Sprint(args...))
}

// Genera un item e lo invia al writer di default del rispettivo livello.
func (w *Writer) Feedbackf(level logs.Level, format string, args ...any) {
	if w.feedbackCh == nil {
		return
	}

	w.feedbackCh <- logs.NewItem(level, w.ID(), fmt.Sprintf(format, args...))
}

// Incapsula e invia un errore al writer di default del livello ErrorLevel.
func (w *Writer) FeedbackError(err error) {
	if w.feedbackCh == nil {
		return
	}

	w.feedbackCh <- logs.NewErrorItem(err)
}
