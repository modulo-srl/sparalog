package logs

// Default logger methods wrapper.

import (
	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logger"
)

// NewLogger allocates, registers and returns a new Logger.
func NewLogger() *logger.Logger {
	if closed {
		return nil
	}

	l := logger.New(DefaultStdoutWriter)

	if l != nil {
		loggers = append(loggers, l)
	}

	return l
}

// NewAliasLogger allocate a logger that uses logs.Default writers.
func NewAliasLogger() sparalog.Logger {
	if closed {
		return nil
	}

	return logger.NewAlias(Default, DefaultDispatcher)
}

// Fatalf sends to Default logger fatal stream using the same fmt.Printf() interface.
func Fatalf(format string, args ...interface{}) {
	Default.Logf(sparalog.FatalLevel, "", format, args...)
}

// Fatal sends to Default logger fatal stream using the same fmt.Print() interface.
func Fatal(args ...interface{}) {
	Default.Log(sparalog.FatalLevel, "", args...)
}

// Errorf sends to Default logger error stream using the same fmt.Printf() interface.
func Errorf(format string, args ...interface{}) {
	Default.Logf(sparalog.ErrorLevel, "", format, args...)
}

// Error sends to Default logger error stream using the same fmt.Print() interface.
func Error(args ...interface{}) {
	Default.Log(sparalog.ErrorLevel, "", args...)
}

// Warnf sends to Default logger warning stream using the same fmt.Printf() interface.
func Warnf(format string, args ...interface{}) {
	Default.Logf(sparalog.WarnLevel, "", format, args...)
}

// Warn sends to Default logger warning stream using the same fmt.Print() interface.
func Warn(args ...interface{}) {
	Default.Log(sparalog.WarnLevel, "", args...)
}

// Infof sends to Default logger info stream using the same fmt.Printf() interface.
func Infof(format string, args ...interface{}) {
	Default.Logf(sparalog.InfoLevel, "", format, args...)
}

// Info sends to Default logger info stream using the same fmt.Print() interface.
func Info(args ...interface{}) {
	Default.Log(sparalog.InfoLevel, "", args...)
}

// Debugf sends to Default logger debug stream using the same fmt.Printf() interface.
func Debugf(format string, args ...interface{}) {
	Default.Logf(sparalog.DebugLevel, "", format, args...)
}

// Debug sends to Default logger debug stream using the same fmt.Print() interface.
func Debug(args ...interface{}) {
	Default.Log(sparalog.DebugLevel, "", args...)
}

// Tracef sends to Default logger trace stream using the same fmt.Printf() interface.
func Tracef(format string, args ...interface{}) {
	Default.Logf(sparalog.TraceLevel, "", format, args...)
}

// Trace sends to Default logger trace stream using the same fmt.Print() interface.
func Trace(args ...interface{}) {
	Default.Log(sparalog.TraceLevel, "", args...)
}

// Printf sends to Default logger info stream using the same fmt.Printf() interface.
func Printf(format string, args ...interface{}) {
	Default.Logf(sparalog.InfoLevel, "", format, args...)
}

// Print sends to Default logger info stream using the same fmt.Print() interface.
func Print(args ...interface{}) {
	Default.Log(sparalog.InfoLevel, "", args...)
}

// ResetWriters reset the writers for all the levels to an optional default writer.
func ResetWriters(defaultW sparalog.Writer) {
	DefaultDispatcher.ResetWriters(defaultW)
}

// ResetLevelWriters remove all level's writers and reset to an optional default writer.
func ResetLevelWriters(level sparalog.Level, defaultW sparalog.Writer) {
	DefaultDispatcher.ResetLevelWriters(level, defaultW)
}

// ResetLevelsWriters remove specific levels writers and reset to an optional default writer.
func ResetLevelsWriters(levels []sparalog.Level, defaultW sparalog.Writer) {
	DefaultDispatcher.ResetLevelsWriters(levels, defaultW)
}

// AddWriter add a writer to all levels.
// id is optional, but useful for RemoveWriter().
func AddWriter(w sparalog.Writer, id sparalog.WriterID) {
	DefaultDispatcher.AddWriter(w, id)
}

// AddLevelWriter add a writer to a level.
// id is optional, but useful for RemoveWriter().
func AddLevelWriter(level sparalog.Level, w sparalog.Writer, id sparalog.WriterID) {
	DefaultDispatcher.AddLevelWriter(level, w, id)
}

// AddLevelsWriter add a writer to several levels.
// id is optional, but useful for RemoveWriter().
func AddLevelsWriter(levels []sparalog.Level, w sparalog.Writer, id sparalog.WriterID) {
	DefaultDispatcher.AddLevelsWriter(levels, w, id)
}

// RemoveWriter delete a specific writer from level.
func RemoveWriter(level sparalog.Level, id sparalog.WriterID) {
	DefaultDispatcher.RemoveWriter(level, id)
}

// Mute mute/unmute a specific level.
func Mute(level sparalog.Level, state bool) {
	DefaultDispatcher.Mute(level, state)
}

// EnableStacktrace enable stacktrace for a specific level.
func EnableStacktrace(level sparalog.Level, state bool) {
	DefaultDispatcher.EnableStacktrace(level, state)
}

// SetContextTag sets a context tag.
func SetContextTag(name, value string) {
	Default.SetContextTag(name, value)
}

// SetContextData sets a context data payload.
func SetContextData(key string, value interface{}) {
	Default.SetContextData(key, value)
}

// SetContextPrefix sets the context prefix.
// Tags is the list of the tags names that will be rendered according to format.
func SetContextPrefix(format string, tags []string) {
	Default.SetContextPrefix(format, tags)
}
