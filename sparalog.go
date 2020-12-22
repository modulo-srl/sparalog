package sparalog

// This is the model.
// For the domain methods see `logs` package.

// Logger is the base interface of the logger.
type Logger interface {
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Tracef(format string, args ...interface{})
	Printf(format string, args ...interface{}) // Infof() alias

	Fatal(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Trace(args ...interface{})
	Print(args ...interface{}) // Info() alias

	SetContextTag(string, string)
	SetContextData(string, interface{})
	SetContextPrefix(format string, tags []string)

	NewItem(Level, ...string) Item
	NewItemf(Level, string, ...string) Item
	NewError(error) Item
	NewErrorf(string, ...interface{}) Item

	// Internal domain functions - should not be used out of this library.

	Log(Level, string, ...interface{})
	Logf(Level, string, string, ...interface{})

	CloneContext() Context

	Close()
}

// Dispatcher is the base class of the logger dispatcher.
type Dispatcher interface {
	Mute(Level, bool)

	EnableStacktrace(Level, bool)

	ResetWriters(Writer)
	ResetLevelWriters(Level, Writer)
	ResetLevelsWriters([]Level, Writer)
	AddWriter(Writer, WriterID)
	AddLevelWriter(Level, Writer, WriterID)
	AddLevelsWriter([]Level, Writer, WriterID)
	RemoveWriter(Level, WriterID)

	// Internal domain functions - should not be used out of this library.

	LevelState(Level) LevelState

	Dispatch(Item)

	Close()
}

// LevelState is the status of specific logging level, returned by Dispatcher.LevelState().
type LevelState struct {
	Muted      bool
	Stacktrace bool
}

// Writer is the writer used by the Logger for one or more log levels.
type Writer interface {
	SetFeedbackChan(chan Item)

	Feedback(Level, ...interface{})
	Feedbackf(Level, string, ...interface{})
	FeedbackItem(Item)

	Write(Item)

	Close()
}

// WriterID is the Writer identifier in the level writers list.
type WriterID string

// FatalExitCode is the Exit Code used in Fatal() and Fatalf().
var FatalExitCode = 1
