package sparalog

// This is the model.
// For the domain methods see `logs` package.

// Logger is the base interface of the logger
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

	ResetWriters(Writer)
	ResetLevelWriters(Level, Writer)
	ResetLevelsWriters([]Level, Writer)
	AddWriter(Writer, WriterID)
	AddLevelWriter(Level, Writer, WriterID)
	AddLevelsWriter([]Level, Writer, WriterID)
	RemoveWriter(Level, WriterID)

	Mute(Level, bool)
	EnableStacktrace(Level, bool)

	// Internal domain functions - should not be used out of this library.

	FatalTrace(stackTrace string, args ...interface{})

	Log(Level, string, bool, ...interface{})
	Logf(Level, string, bool, string, ...interface{})

	Write(Item, bool)

	Close()
}

// Writer is the writer used by the Logger for one or more log levels.
type Writer interface {
	Write(Item) WriterError
	Close()

	SetFeedbackChan(chan WriterError)
	FeedbackError(WriterError)

	ErrorItem(error) WriterError
}

// WriterID is the Writer ID.
type WriterID string

// WriterError is an error, with stack trace, wrapped on Item.
type WriterError *Item

// FatalExitCode is the Exit Code used in Fatal() and Fatalf().
var FatalExitCode = 1
