package logger

import (
	"fmt"
	"sync"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/item"
)

// New allocate a new logger.
func New(defaultWriter sparalog.Writer) *Logger {
	l := Logger{
		context: &item.Context{},
	}

	l.initDispatcher(defaultWriter)

	return &l
}

// NewAlias allocate a logger using another logger context and another dispatcher.
func NewAlias(logger sparalog.Logger, dispatcher sparalog.Dispatcher) sparalog.Logger {
	if logger == nil || dispatcher == nil {
		return nil
	}

	l := Logger{
		context:    logger.CloneContext(),
		Dispatcher: dispatcher,
	}

	return &l
}

// Logger implements Logger inferface.
type Logger struct {
	Dispatcher sparalog.Dispatcher

	muCtx   sync.RWMutex
	context sparalog.Context
}

// Fatalf logs to fatal stream using the same fmt.Printf() interface.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logf(sparalog.FatalLevel, "", format, args...)
}

// Fatal logs to fatal stream using the same fmt.Print() interface.
func (l *Logger) Fatal(args ...interface{}) {
	l.Log(sparalog.FatalLevel, "", args...)
}

// Errorf logs to error stream using the same fmt.Printf() interface.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logf(sparalog.ErrorLevel, "", format, args...)
}

// Error logs to error stream using the same fmt.Print() interface.
func (l *Logger) Error(args ...interface{}) {
	l.Log(sparalog.ErrorLevel, "", args...)
}

// Warnf logs to warning stream using the same fmt.Printf() interface.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logf(sparalog.WarnLevel, "", format, args...)
}

// Warn logs to warning stream using the same fmt.Print() interface.
func (l *Logger) Warn(args ...interface{}) {
	l.Log(sparalog.WarnLevel, "", args...)
}

// Infof logs to info stream using the same fmt.Printf() interface.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logf(sparalog.InfoLevel, "", format, args...)
}

// Info logs to info stream using the same fmt.Print() interface.
func (l *Logger) Info(args ...interface{}) {
	l.Log(sparalog.InfoLevel, "", args...)
}

// Printf logs to info stream using the same fmt.Printf() interface.
func (l *Logger) Printf(format string, args ...interface{}) {
	l.Logf(sparalog.InfoLevel, "", format, args...)
}

// Print logs to info stream using the same fmt.Print() interface.
func (l *Logger) Print(args ...interface{}) {
	l.Log(sparalog.InfoLevel, "", args...)
}

// Debugf logs to debug stream using the same fmt.Printf() interface.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logf(sparalog.DebugLevel, "", format, args...)
}

// Debug logs to debug stream using the same fmt.Print() interface.
func (l *Logger) Debug(args ...interface{}) {
	l.Log(sparalog.DebugLevel, "", args...)
}

// Tracef logs to trace stream using the same fmt.Printf() interface.
func (l *Logger) Tracef(format string, args ...interface{}) {
	l.Logf(sparalog.TraceLevel, "", format, args...)
}

// Trace logs to trace stream using the same fmt.Print() interface.
func (l *Logger) Trace(args ...interface{}) {
	l.Log(sparalog.TraceLevel, "", args...)
}

// NewItem generate a new log item.
func (l *Logger) NewItem(level sparalog.Level, args ...string) sparalog.Item {
	i := item.New(level, fmt.Sprint(args))
	i.SetLogger(l)

	l.muCtx.RLock()
	defer l.muCtx.RUnlock()
	i.AssignContext(l.context, true)

	return i
}

// NewItemf generate a new log item.
func (l *Logger) NewItemf(level sparalog.Level, format string, args ...string) sparalog.Item {
	i := item.New(level, fmt.Sprintf(format, args))
	i.SetLogger(l)

	l.muCtx.RLock()
	defer l.muCtx.RUnlock()
	i.AssignContext(l.context, true)

	return i
}

// NewErrorItem generate a new log error item with stack trace, excluding the caller.
// Entry point for all the NewError* helpers.
func (l *Logger) NewErrorItem(line string) sparalog.Item {
	state := l.Dispatcher.LevelState(sparalog.ErrorLevel)

	i := item.New(sparalog.ErrorLevel, line)

	if state.Stacktrace {
		i.GenerateStackTrace(2) // current + caller
	}

	i.SetLogger(l)

	l.muCtx.RLock()
	defer l.muCtx.RUnlock()
	i.AssignContext(l.context, true)

	return i
}

// NewError generate a new log item wrapping an error.
func (l *Logger) NewError(err error) sparalog.Item {
	line := err.Error() // as item.NewError()
	return l.NewErrorItem(line)
}

// NewErrorf generate a new log item wrapping an error, as Errorf().
func (l *Logger) NewErrorf(format string, a ...interface{}) sparalog.Item {
	line := fmt.Errorf(format, a...).Error() // as item.NewErrorf()
	return l.NewErrorItem(line)
}

// SetContextTag sets a context tag.
func (l *Logger) SetContextTag(name, value string) {
	l.muCtx.Lock()
	defer l.muCtx.Unlock()

	l.context.SetTag(name, value)
}

// SetContextData sets a context data payload.
func (l *Logger) SetContextData(key string, value interface{}) {
	l.muCtx.Lock()
	defer l.muCtx.Unlock()

	l.context.SetData(key, value)
}

// SetContextPrefix sets the context prefix.
// Tags is the list of the tags names that will be rendered according to format.
func (l *Logger) SetContextPrefix(format string, tags []string) {
	l.muCtx.Lock()
	defer l.muCtx.Unlock()

	l.context.SetPrefix(format, tags)
}

// prelogItem generate a prefilled log item - thread safe.
// Returns nil when the log cannot be performed (level is muted).
// stackTrace: custom stacktrace; "" = generate here according to level settings.
func (l *Logger) prelogItem(level sparalog.Level, stackTrace string) sparalog.Item {
	state := l.Dispatcher.LevelState(level)
	if state.Muted || state.NoWriters {
		return nil
	}

	i := item.New(level, "")

	if state.Stacktrace && stackTrace == "" {
		i.GenerateStackTrace(3) // this call + Log*() + level-helper()
	}

	l.muCtx.RLock()
	defer l.muCtx.RUnlock()
	i.AssignContext(l.context, false)

	return i
}

// Log to level stream - entry point for all the helpers; thread safe.
func (l *Logger) Log(level sparalog.Level, stackTrace string, args ...interface{}) {
	item := l.prelogItem(level, stackTrace)
	if item == nil {
		return
	}

	item.SetLine(fmt.Sprint(args...))

	l.Dispatcher.Dispatch(item)
}

// Logf logs to level stream using format - entry point for all the helpers; thread safe.
func (l *Logger) Logf(level sparalog.Level, stackTrace string, format string, args ...interface{}) {
	item := l.prelogItem(level, stackTrace)
	if item == nil {
		return
	}

	item.SetLine(fmt.Sprintf(format, args...))

	l.Dispatcher.Dispatch(item)
}

// LogItem logs an item to level stream - thread safe.
func (l *Logger) LogItem(item sparalog.Item) {
	state := l.Dispatcher.LevelState(item.Level())
	if state.Muted || state.NoWriters {
		return
	}

	if !state.Stacktrace {
		item.SetStackTrace("")
	}

	l.Dispatcher.Dispatch(item)
}

// CloneContext returns a clone of the current context.
func (l *Logger) CloneContext() sparalog.Context {
	l.muCtx.RLock()
	defer l.muCtx.RUnlock()

	ctx := &item.Context{}
	ctx.AssignContext(l.context, true)

	return ctx
}

func (l *Logger) initDispatcher(defaultWriter sparalog.Writer) {
	l.Dispatcher = NewDispatcher(defaultWriter)

	l.Dispatcher.EnableStacktrace(sparalog.FatalLevel, true)
	l.Dispatcher.EnableStacktrace(sparalog.ErrorLevel, true)

	l.Dispatcher.Mute(sparalog.DebugLevel, true)
	l.Dispatcher.Mute(sparalog.TraceLevel, true)
}

// Open start loggers and all the writers.
func (l *Logger) Open() error {
	return l.Dispatcher.Start()
}

// Close terminates loggers and all the writers.
func (l *Logger) Close() {
	l.Dispatcher.Stop()
}

var defaultWriterID sparalog.WriterID = "0"
