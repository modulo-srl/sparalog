package sparalog

// This is the model.
// For the domain methods see `logs` package.

// Logger is the base interface of the logger.
type Logger interface {
	// Fatal sends to Default logger fatal stream using the same fmt.Print() interface.
	Fatal(args ...interface{})
	// Fatalf sends to Default logger fatal stream using the same fmt.Printf() interface.
	Fatalf(format string, args ...interface{})

	// Error sends to Default logger error stream using the same fmt.Print() interface.
	Error(args ...interface{})
	// Errorf sends to Default logger error stream using the same fmt.Printf() interface.
	Errorf(format string, args ...interface{})

	// Warn sends to Default logger warning stream using the same fmt.Print() interface.
	Warn(args ...interface{})
	// Warnf sends to Default logger warning stream using the same fmt.Printf() interface.
	Warnf(format string, args ...interface{})

	// Info sends to Default logger info stream using the same fmt.Print() interface.
	Info(args ...interface{})
	// Infof sends to Default logger info stream using the same fmt.Printf() interface.
	Infof(format string, args ...interface{})

	// Debug sends to Default logger debug stream using the same fmt.Print() interface.
	Debug(args ...interface{})
	// Debugf sends to Default logger debug stream using the same fmt.Printf() interface.
	Debugf(format string, args ...interface{})

	// Trace sends to Default logger trace stream using the same fmt.Print() interface.
	Trace(args ...interface{})
	// Tracef sends to Default logger trace stream using the same fmt.Printf() interface.
	Tracef(format string, args ...interface{})

	// Print sends to Default logger info stream using the same fmt.Print() interface.
	Print(args ...interface{}) // Info() alias
	// Printf sends to Default logger info stream using the same fmt.Printf() interface.
	Printf(format string, args ...interface{}) // Infof() alias

	// SetContextTag sets a context tag.
	SetContextTag(string, string)
	// SetContextData sets a context data payload.
	SetContextData(string, interface{})
	// SetContextPrefix sets the context prefix.
	// Tags is the list of the tags names that will be rendered according to format.
	SetContextPrefix(format string, tags []string)

	// NewItem generate a new log item for the default logger.
	NewItem(Level, ...string) Item
	// NewItemf generate a new log item for the default logger.
	NewItemf(Level, string, ...string) Item
	// NewError generate a new log item wrapping an error.
	NewError(error) Item
	// NewErrorf generate a new log item wrapping an error, as Errorf().
	NewErrorf(string, ...interface{}) Item

	// Internal domain functions - should not be used out of this library.

	// Log to level stream - entry point for all the helpers; thread safe.
	Log(Level, string, ...interface{})
	// Logf logs to level stream using format - entry point for all the helpers; thread safe.
	Logf(Level, string, string, ...interface{})

	// CloneContext returns a clone of the current context.
	CloneContext() Context

	// Close terminates loggers and all the writers.
	Close()
}

// Dispatcher is the base class of the logger dispatcher.
type Dispatcher interface {
	// Mute mute/unmute a specific level.
	Mute(Level, bool)

	// EnableStacktrace enable stacktrace for a specific level.
	EnableStacktrace(Level, bool)

	// ResetWriters reset the writers for all the levels to an optional default writer.
	ResetWriters(Writer)
	// ResetLevelWriters remove all level's writers and reset to an optional default writer.
	ResetLevelWriters(Level, Writer)
	// ResetLevelsWriters remove specific levels writers and reset to an optional default writer.
	ResetLevelsWriters([]Level, Writer)
	// AddWriter add a writer to all levels.
	AddWriter(Writer, WriterID)
	// AddLevelWriter add a writer to a level.
	AddLevelWriter(Level, Writer, WriterID)
	// AddLevelsWriter add a writer to several levels.
	AddLevelsWriter([]Level, Writer, WriterID)
	// RemoveWriter delete a specific writer from level.
	RemoveWriter(Level, WriterID)

	// Internal domain functions - should not be used out of this library.

	// LevelState return a level status.
	LevelState(Level) LevelState

	// Dispatch sends an Item to the level writers.
	Dispatch(Item)

	// Close terminate all the writers.
	Close()
}

// LevelState is the status of specific logging level, returned by Dispatcher.LevelState().
type LevelState struct {
	Muted      bool
	Stacktrace bool
}

// Writer is the writer used by the Logger for one or more log levels.
type Writer interface {
	// SetFeedbackChan set a channel to the level default writer of the logger.
	SetFeedbackChan(chan Item)

	// Feedback generate an item and send it to the level default writer of the logger.
	Feedback(Level, ...interface{})
	// Feedbackf generate an item and send it to the level default writer of the logger.
	Feedbackf(Level, string, ...interface{})
	// FeedbackItem send an item to the level default writer of ther logger.
	FeedbackItem(Item)

	// Write writes an item.
	Write(Item)

	// Close terminates the writer.
	Close()
}

// WriterID is the Writer identifier in the level writers list.
type WriterID string

// FatalExitCode is the Exit Code used in Fatal() and Fatalf().
var FatalExitCode = 1
