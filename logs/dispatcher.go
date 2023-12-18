package logs

// Dispatcher del logger.

import (
	"os"
	"runtime"
	"sync"
	"time"
)

type dispatcher struct {
	levelWriters [LevelsCount]levelWriters

	writersFeedback   chan *Item
	writersFeedbackWG sync.WaitGroup

	mu     sync.RWMutex
	muted  [LevelsCount]bool
	closed bool
}

type levelWriters struct {
	writers       map[Writer]bool
	defaultWriter Writer
}

// Alloca e inizializza un nuovo dispatcher.
func newDispatcher(defaultWriter Writer) *dispatcher {
	d := dispatcher{
		writersFeedback: make(chan *Item, 64),
	}

	d.ResetWriters(defaultWriter)

	d.startFeedbackWatcher()

	// For non default loggers only, because is not called after main termination.
	runtime.SetFinalizer(&d, finalizeDispatcher)

	return &d
}

// Disassocia tutti i writer e reimposta un writer di default
// (il writer di default riceve le eventuali loggate di errore o di feedback da parte degli altri writer).
// NON thread safe.
func (d *dispatcher) ResetWriters(defaultW Writer) {
	for i := 0; i < int(LevelsCount); i++ {
		d.ResetLevelWriters(Level(i), defaultW)
	}
}

// Disassocia tutti i writer per un certo livello e ne reimposta un writer di default
// (il writer di default riceve le eventuali loggate di errore o di feedback da parte degli altri writer).
// NON thread safe.
func (d *dispatcher) ResetLevelWriters(level Level, defaultW Writer) {
	l := levelWriters{
		writers:       make(map[Writer]bool),
		defaultWriter: defaultW,
	}

	if defaultW != nil {
		l.writers[defaultW] = true
	}

	d.levelWriters[level] = l
}

// Disassocia tutti i writer per un set di livelli e ne reimposta un writer di default
// (il writer di default riceve le eventuali loggate di errore o di feedback da parte degli altri writer).
// NON thread safe.
func (d *dispatcher) ResetLevelsWriters(levels []Level, defaultW Writer) {
	for _, level := range levels {
		d.ResetLevelWriters(level, defaultW)
	}
}

// Associa un writer a tutti i livelli.
// NON thread safe.
func (d *dispatcher) AddWriter(w Writer) {
	for level := 0; level < int(LevelsCount); level++ {
		d.levelWriters[level].writers[w] = true
	}

	w.SetFeedbackChan(d.writersFeedback)
}

// Associa un writer a uno specifico livello.
// NON thread safe.
func (d *dispatcher) AddLevelWriter(level Level, w Writer) {
	d.levelWriters[level].writers[w] = true

	w.SetFeedbackChan(d.writersFeedback)
}

// Associa un writer a un set di livelli.
// NON thread safe.
func (d *dispatcher) AddLevelsWriter(levels []Level, w Writer) {
	for _, level := range levels {
		d.levelWriters[level].writers[w] = true
	}

	w.SetFeedbackChan(d.writersFeedback)
}

// Muta o smuta un livello.
func (d *dispatcher) Mute(level Level, state bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.muted[level] = state
}

// Invia un item a tutti i writer del livello.
func (d *dispatcher) Dispatch(item *Item) {
	for w := range d.levelWriters[item.Level].writers {
		w.Write(item)
	}

	if item.Level == FatalLevel {
		d.Stop()
		os.Exit(FatalExitCode)
	}
}

// Avvia tutti i writer.
func (d *dispatcher) Start() error {
	openedWw := make(map[Writer]bool)

	for _, lw := range d.levelWriters {
		for w := range lw.writers {
			// Assicura una sola chiamata a writer.Open()
			if _, ok := openedWw[w]; ok {
				continue
			}

			err := w.Start()
			if err != nil {
				return err
			}

			openedWw[w] = true
		}
	}

	return nil
}

// Stoppa tutti i writer e il canale di feedback.
func (d *dispatcher) Stop() {
	d.mu.Lock()
	if d.closed {
		d.mu.Unlock()
		return
	}
	d.closed = true
	d.mu.Unlock()

	// Chiude e vuota il canale di feedback.
	close(d.writersFeedback)
	waitTimeout(&d.writersFeedbackWG, time.Second*time.Duration(3))

	// Stoppa i writer attendendo che terminino le proprie code.
	closedWw := make(map[Writer]bool)

	for _, lw := range d.levelWriters {
		for w := range lw.writers {
			// Assicura una singola chiamata a writer.Stop()
			if _, ok := closedWw[w]; ok {
				continue
			}

			w.Stop()
			closedWw[w] = true
		}
	}
}

// Ritorna true se il livello ha almeno un writer che non sia mutato.
func (d *dispatcher) CanDispatch(level Level) bool {
	if len(d.levelWriters[level].writers) == 0 {
		return false
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.closed || d.muted[level] {
		return false
	}

	return true
}

func (d *dispatcher) startFeedbackWatcher() {
	d.writersFeedbackWG.Add(1)

	go func() {
		for item := range d.writersFeedback {
			w := d.levelWriters[item.Level].defaultWriter
			if w != nil {
				w.Write(item)
			}

			if item.Level == FatalLevel {
				d.Stop()
				os.Exit(FatalExitCode)
			}
		}

		d.writersFeedbackWG.Done()
	}()
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

func finalizeDispatcher(d *dispatcher) {
	d.Stop()
}
