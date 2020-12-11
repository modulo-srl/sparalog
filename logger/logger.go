package logger

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/env"
)

// New allocate a new logger.
// Tag is used as logging prefix.
func New(tag string, defaultWriter sparalog.Writer) sparalog.Logger {
	l := Logger{}

	l.init(tag, defaultWriter)

	return &l
}

// NewChild allocate a child logger that uses parentLogger writers.
// If parent = nil the logger will be a child of the logger.Default.
func NewChild(parent sparalog.Logger, tag string) sparalog.Logger {
	if parent == nil {
		return nil
	}

	l := Logger{
		tag:    tag,
		parent: parent,
	}

	l.init(tag, nil)

	return &l
}

// Logger implements Logger inferface.
type Logger struct {
	// Set to forward to parent.LogString(), using parent's writers.
	parent sparalog.Logger

	tag string

	writersErrorFeedback chan sparalog.WriterError
	writersFeedbackWG    sync.WaitGroup

	mu       sync.RWMutex
	writers  [sparalog.LevelsCount]map[sparalog.WriterID]sparalog.Writer
	levState [sparalog.LevelsCount]levelState

	closed bool
}

type levelState struct {
	muted      bool
	stacktrace bool
}

// Fatalf logs to fatal stream using the same fmt.Printf() interface.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logf(sparalog.FatalLevel, "", false, format, args...)
}

// Fatal logs to fatal stream using the same fmt.Print() interface.
func (l *Logger) Fatal(args ...interface{}) {
	l.Log(sparalog.FatalLevel, "", false, args...)
}

// FatalTrace logs to fatal stream with a custom stack trace.
func (l *Logger) FatalTrace(stackTrace string, args ...interface{}) {
	l.Log(sparalog.FatalLevel, stackTrace, false, args...)
}

// Errorf logs to error stream using the same fmt.Printf() interface.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logf(sparalog.ErrorLevel, "", false, format, args...)
}

// Error logs to error stream using the same fmt.Print() interface.
func (l *Logger) Error(args ...interface{}) {
	l.Log(sparalog.ErrorLevel, "", false, args...)
}

// Warnf logs to warning stream using the same fmt.Printf() interface.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logf(sparalog.WarnLevel, "", false, format, args...)
}

// Warn logs to warning stream using the same fmt.Print() interface.
func (l *Logger) Warn(args ...interface{}) {
	l.Log(sparalog.WarnLevel, "", false, args...)
}

// Infof logs to info stream using the same fmt.Printf() interface.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logf(sparalog.InfoLevel, "", false, format, args...)
}

// Info logs to info stream using the same fmt.Print() interface.
func (l *Logger) Info(args ...interface{}) {
	l.Log(sparalog.InfoLevel, "", false, args...)
}

// Printf logs to info stream using the same fmt.Printf() interface.
func (l *Logger) Printf(format string, args ...interface{}) {
	l.Logf(sparalog.InfoLevel, "", false, format, args...)
}

// Print logs to info stream using the same fmt.Print() interface.
func (l *Logger) Print(args ...interface{}) {
	l.Log(sparalog.InfoLevel, "", false, args...)
}

// Debugf logs to debug stream using the same fmt.Printf() interface.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logf(sparalog.DebugLevel, "", false, format, args...)
}

// Debug logs to debug stream using the same fmt.Print() interface.
func (l *Logger) Debug(args ...interface{}) {
	l.Log(sparalog.DebugLevel, "", false, args...)
}

// Tracef logs to trace stream using the same fmt.Printf() interface.
func (l *Logger) Tracef(format string, args ...interface{}) {
	l.Logf(sparalog.TraceLevel, "", false, format, args...)
}

// Trace logs to trace stream using the same fmt.Print() interface.
func (l *Logger) Trace(args ...interface{}) {
	l.Log(sparalog.TraceLevel, "", false, args...)
}

// ResetWriters reset the writers for all the levels to an optional default writer.
func (l *Logger) ResetWriters(defaultW sparalog.Writer) {
	if l.parent != nil {
		return
	}

	for i := range l.writers {
		l.ResetLevelWriters(sparalog.Level(i), defaultW)
	}
}

// ResetLevelWriters remove all level's writers and reset to an optional default writer.
func (l *Logger) ResetLevelWriters(level sparalog.Level, defaultW sparalog.Writer) {
	if l.parent != nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	ww := make(map[sparalog.WriterID]sparalog.Writer, 4)

	if defaultW != nil {
		ww[defaultWriterID] = defaultW
	}

	l.writers[level] = ww
}

// AddWriter add a writer to all levels.
// id is optional, but useful for RemoveWriter().
func (l *Logger) AddWriter(w sparalog.Writer, id sparalog.WriterID) {
	if l.parent != nil {
		return
	}

	for i := range l.writers {
		l.AddLevelWriter(sparalog.Level(i), w, id)
	}
}

// AddLevelWriter add a writer to a level.
// id is optional, but useful for RemoveWriter().
func (l *Logger) AddLevelWriter(level sparalog.Level, w sparalog.Writer, id sparalog.WriterID) {
	if l.parent != nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if id == "" {
		id = sparalog.WriterID(strconv.Itoa(len(l.writers[level])))
	}

	l.writers[level][id] = w

	if id != defaultWriterID {
		w.SetFeedbackChan(l.writersErrorFeedback)
	}
}

// AddLevelsWriter add a writer to several levels.
// id is optional, but useful for RemoveWriter().
func (l *Logger) AddLevelsWriter(levels []sparalog.Level, w sparalog.Writer, id sparalog.WriterID) {
	if l.parent != nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	for level := range levels {
		if id == "" {
			id = sparalog.WriterID(strconv.Itoa(len(l.writers[level])))
		}

		l.writers[level][id] = w
	}

	if id != defaultWriterID {
		w.SetFeedbackChan(l.writersErrorFeedback)
	}
}

// RemoveWriter delete a specific writer from level.
func (l *Logger) RemoveWriter(level sparalog.Level, id sparalog.WriterID) {
	if l.parent != nil {
		return
	}

	if id == "" {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.writers[level], id)
}

// Mute mute/unmute a specific level.
func (l *Logger) Mute(level sparalog.Level, state bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.levState[level].muted = state
}

// EnableStacktrace enable stacktrace for a specific level.
func (l *Logger) EnableStacktrace(level sparalog.Level, state bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.levState[level].stacktrace = state
}

// newLogItem generate a prefilled log item - thread safe.
// If ok = false the log cannot be performed (level is muted).
// stackTrace: custom stacktrace.
func (l *Logger) newLogItem(level sparalog.Level, stackTrace string) (item sparalog.Item, ok bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.closed || l.levState[level].muted {
		ok = false
		return
	}

	var st string
	if l.levState[level].stacktrace {
		if stackTrace != "" {
			st = stackTrace
		} else {
			st = env.StackTrace(skipStackCalls)
		}
	}

	item = sparalog.Item{
		Timestamp:  time.Now(),
		Level:      level,
		Tag:        l.tag,
		StackTrace: st,
	}

	ok = true
	return
}

// Log to level stream - abstract function, thread safe.
func (l *Logger) Log(level sparalog.Level, stackTrace string, defaultWriterOnly bool, args ...interface{}) {
	defer func() {
		if level == sparalog.FatalLevel {
			l.Close()
			os.Exit(sparalog.FatalExitCode)
		}
	}()

	item, ok := l.newLogItem(level, stackTrace)
	if !ok {
		return
	}

	item.Line = fmt.Sprint(args...)

	if l.parent != nil {
		l.parent.Write(item, defaultWriterOnly)
		return
	}

	l.Write(item, defaultWriterOnly)
}

// Logf logs to level stream using format - abstract function, thread safe.
func (l *Logger) Logf(level sparalog.Level, stackTrace string, defaultWriterOnly bool, format string, args ...interface{}) {
	defer func() {
		if level == sparalog.FatalLevel {
			l.Close()
			os.Exit(sparalog.FatalExitCode)
		}
	}()

	item, ok := l.newLogItem(level, stackTrace)
	if !ok {
		return
	}

	item.Line = fmt.Sprintf(format, args...)

	if l.parent != nil {
		l.parent.Write(item, defaultWriterOnly)
		return
	}

	l.Write(item, defaultWriterOnly)
}

// Write sends an Item to the level writers.
func (l *Logger) Write(item sparalog.Item, defaultWriterOnly bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if defaultWriterOnly {
		w, ok := l.writers[item.Level][defaultWriterID]
		if ok {
			w.Write(item)
		}

		return
	}

	for id, w := range l.writers[item.Level] {
		if w == nil {
			continue
		}

		err := w.Write(item)

		if err != nil && id != defaultWriterID {
			l.writersErrorFeedback <- err
		}
	}
}

func (l *Logger) init(tag string, defaultWriter sparalog.Writer) {
	l.tag = tag

	l.writersErrorFeedback = make(chan sparalog.WriterError, 64)

	l.EnableStacktrace(sparalog.FatalLevel, true)

	if l.parent == nil {
		l.ResetWriters(defaultWriter)
	}

	l.startErrorFeedbackWatcher()

	// For non default loggers only, because is not called after main termination.
	runtime.SetFinalizer(l, finalizeLogger)
}

func (l *Logger) startErrorFeedbackWatcher() {
	l.writersFeedbackWG.Add(1)

	go func() {
		for item := range l.writersErrorFeedback {
			l.Log(item.Level, item.StackTrace, true, item.Line)
		}

		l.writersFeedbackWG.Done()
	}()
}

// Close terminate loggers and all the writers.
func (l *Logger) Close() {
	l.mu.Lock()

	if l.closed {
		l.mu.Unlock()
		return
	}
	l.closed = true
	l.mu.Unlock()

	// Close and wait for writers feedback channel.
	close(l.writersErrorFeedback)
	waitTimeout(&l.writersFeedbackWG, time.Second*time.Duration(3))

	l.mu.Lock()
	defer l.mu.Unlock()

	// Close and wait for writers.
	closedWw := make(map[sparalog.Writer]bool)

	for level := range l.writers {
		for _, w := range l.writers[level] {
			// Assure single call of writer.Close()
			if _, ok := closedWw[w]; ok {
				continue
			}

			w.Close()
			closedWw[w] = true
		}
	}
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

func finalizeLogger(l *Logger) {
	l.Close()
}

// How many top calls to skip from the stack trace.
var skipStackCalls = 5

var defaultWriterID sparalog.WriterID = "0"
