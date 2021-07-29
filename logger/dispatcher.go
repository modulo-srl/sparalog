package logger

import (
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/modulo-srl/sparalog"
)

// NewDispatcher allocate and initialize a new dispatcher.
func NewDispatcher(defaultWriter sparalog.Writer) *Dispatcher {
	l := Dispatcher{}

	l.init(defaultWriter)

	return &l
}

// Dispatcher implements sparalog.Dispatcher
type Dispatcher struct {
	writersFeedback   chan sparalog.Item
	writersFeedbackWG sync.WaitGroup

	mu       sync.RWMutex
	writers  [sparalog.LevelsCount]map[sparalog.WriterID]sparalog.Writer
	levState [sparalog.LevelsCount]sparalog.LevelState

	closed bool
}

// ResetWriters reset the writers for all the levels to an optional default writer.
func (d *Dispatcher) ResetWriters(defaultW sparalog.Writer) {
	for i := 0; i < int(sparalog.LevelsCount); i++ {
		d.ResetLevelWriters(sparalog.Level(i), defaultW)
	}
}

// ResetLevelWriters remove all level's writers and reset to an optional default writer.
func (d *Dispatcher) ResetLevelWriters(level sparalog.Level, defaultW sparalog.Writer) {
	d.mu.Lock()
	defer d.mu.Unlock()

	ww := make(map[sparalog.WriterID]sparalog.Writer)

	if defaultW != nil {
		ww[defaultWriterID] = defaultW
	}

	d.writers[level] = ww

	d.levState[level].NoWriters = (defaultW == nil)
}

// ResetLevelsWriters remove specific levels writers and reset to an optional default writer.
func (d *Dispatcher) ResetLevelsWriters(levels []sparalog.Level, defaultW sparalog.Writer) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, level := range levels {
		ww := make(map[sparalog.WriterID]sparalog.Writer)

		if defaultW != nil {
			ww[defaultWriterID] = defaultW
		}

		d.writers[level] = ww

		d.levState[level].NoWriters = (defaultW == nil)
	}
}

// AddWriter add a writer to all levels.
// id is optional, but useful for RemoveWriter().
func (d *Dispatcher) AddWriter(w sparalog.Writer, id sparalog.WriterID) {
	for level := 0; level < int(sparalog.LevelsCount); level++ {
		d.AddLevelWriter(sparalog.Level(level), w, id)
		d.levState[level].NoWriters = false
	}
}

// AddLevelWriter add a writer to a level.
// id is optional, but useful for RemoveWriter().
func (d *Dispatcher) AddLevelWriter(level sparalog.Level, w sparalog.Writer, id sparalog.WriterID) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if id == "" {
		id = sparalog.WriterID(strconv.Itoa(len(d.writers[level])))
	}

	d.writers[level][id] = w

	d.levState[level].NoWriters = false

	if id != defaultWriterID {
		w.SetFeedbackChan(d.writersFeedback)
	}
}

// AddLevelsWriter add a writer to several levels.
// id is optional, but useful for RemoveWriter().
func (d *Dispatcher) AddLevelsWriter(levels []sparalog.Level, w sparalog.Writer, id sparalog.WriterID) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, level := range levels {
		if id == "" {
			id = sparalog.WriterID(strconv.Itoa(len(d.writers[level])))
		}

		d.writers[level][id] = w

		d.levState[level].NoWriters = false
	}

	if id != defaultWriterID {
		w.SetFeedbackChan(d.writersFeedback)
	}
}

// RemoveWriter delete a specific writer from level.
func (d *Dispatcher) RemoveWriter(level sparalog.Level, id sparalog.WriterID) {
	if id == "" {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.writers[level], id)

	d.levState[level].NoWriters = true
	for _, w := range d.writers[level] {
		if w != nil {
			d.levState[level].NoWriters = false
			break
		}
	}
}

// Mute mute/unmute a specific level.
func (d *Dispatcher) Mute(level sparalog.Level, state bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.levState[level].Muted = state
}

// EnableStacktrace enable stacktrace for a specific level.
func (d *Dispatcher) EnableStacktrace(level sparalog.Level, state bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.levState[level].Stacktrace = state
}

// Dispatch sends an Item to the level writers.
func (d *Dispatcher) Dispatch(item sparalog.Item) {
	d.write(item, false)
}

// write writes an item using the level writers or the default writer.
func (d *Dispatcher) write(item sparalog.Item, defaultWriterOnly bool) {
	level := item.Level()

	defer func() {
		if level == sparalog.FatalLevel {
			d.Close()
			os.Exit(sparalog.FatalExitCode)
		}
	}()

	d.mu.RLock()
	defer d.mu.RUnlock()

	if defaultWriterOnly {
		// To default only.
		w, ok := d.writers[level][defaultWriterID]
		if ok {
			w.Write(item)
		}

		return
	}

	// To all writers.
	for _, w := range d.writers[level] {
		if w == nil {
			continue
		}

		w.Write(item)
	}
}

func (d *Dispatcher) init(defaultWriter sparalog.Writer) {
	d.writersFeedback = make(chan sparalog.Item, 64)

	d.ResetWriters(defaultWriter)

	d.startFeedbackWatcher()

	// For non default loggers only, because is not called after main termination.
	runtime.SetFinalizer(d, finalizeDispatcher)
}

func (d *Dispatcher) startFeedbackWatcher() {
	d.writersFeedbackWG.Add(1)

	go func() {
		for item := range d.writersFeedback {
			d.write(item, true)
		}

		d.writersFeedbackWG.Done()
	}()
}

// Close terminate all the writers.
func (d *Dispatcher) Close() {
	d.mu.Lock()

	if d.closed {
		d.mu.Unlock()
		return
	}
	d.closed = true
	d.mu.Unlock()

	// Close and wait for writers feedback channel.
	close(d.writersFeedback)
	waitTimeout(&d.writersFeedbackWG, time.Second*time.Duration(3))

	d.mu.Lock()
	defer d.mu.Unlock()

	// Close and wait for writers.
	closedWw := make(map[sparalog.Writer]bool)

	for level := range d.writers {
		for _, w := range d.writers[level] {
			// Assure single call of writer.Close()
			if _, ok := closedWw[w]; ok {
				continue
			}

			w.Close()
			closedWw[w] = true
		}
	}
}

// LevelState return a level status.
func (d *Dispatcher) LevelState(level sparalog.Level) sparalog.LevelState {
	d.mu.RLock()
	defer d.mu.RUnlock()

	state := d.levState[level]
	if d.closed {
		state.Muted = true
	}

	return state
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

func finalizeDispatcher(d *Dispatcher) {
	d.Close()
}
