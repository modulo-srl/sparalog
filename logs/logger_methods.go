package logs

import (
	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logger"
)

// NewLogger allocates, registers and returns a new Logger.
// Tag is used as logging prefix.
func NewLogger(tag string) sparalog.Logger {
	if closed {
		return nil
	}

	l := logger.New(tag, DefaultStdoutWriter)

	if l != nil {
		loggers = append(loggers, l)
	}

	return l
}

// NewChildLogger allocate a child logger that uses logs.Default writers.
func NewChildLogger(tag string) sparalog.Logger {
	if closed {
		return nil
	}

	return logger.NewChild(Default, tag)
}

// Fatalf sends to Default logger fatal stream using the same fmt.Printf() interface.
func Fatalf(format string, args ...interface{}) {
	Default.Logf(sparalog.FatalLevel, "", false, format, args...)
}

// Fatal sends to Default logger fatal stream using the same fmt.Print() interface.
func Fatal(args ...interface{}) {
	Default.Log(sparalog.FatalLevel, "", false, args...)
}

// Errorf sends to Default logger error stream using the same fmt.Printf() interface.
func Errorf(format string, args ...interface{}) {
	Default.Logf(sparalog.ErrorLevel, "", false, format, args...)
}

// Error sends to Default logger error stream using the same fmt.Print() interface.
func Error(args ...interface{}) {
	Default.Log(sparalog.ErrorLevel, "", false, args...)
}

// Warnf sends to Default logger warning stream using the same fmt.Printf() interface.
func Warnf(format string, args ...interface{}) {
	Default.Logf(sparalog.WarnLevel, "", false, format, args...)
}

// Warn sends to Default logger warning stream using the same fmt.Print() interface.
func Warn(args ...interface{}) {
	Default.Log(sparalog.WarnLevel, "", false, args...)
}

// Infof sends to Default logger info stream using the same fmt.Printf() interface.
func Infof(format string, args ...interface{}) {
	Default.Logf(sparalog.InfoLevel, "", false, format, args...)
}

// Info sends to Default logger info stream using the same fmt.Print() interface.
func Info(args ...interface{}) {
	Default.Log(sparalog.InfoLevel, "", false, args...)
}

// Debugf sends to Default logger debug stream using the same fmt.Printf() interface.
func Debugf(format string, args ...interface{}) {
	Default.Logf(sparalog.DebugLevel, "", false, format, args...)
}

// Debug sends to Default logger debug stream using the same fmt.Print() interface.
func Debug(args ...interface{}) {
	Default.Log(sparalog.DebugLevel, "", false, args...)
}

// Tracef sends to Default logger trace stream using the same fmt.Printf() interface.
func Tracef(format string, args ...interface{}) {
	Default.Logf(sparalog.TraceLevel, "", false, format, args...)
}

// Trace sends to Default logger trace stream using the same fmt.Print() interface.
func Trace(args ...interface{}) {
	Default.Log(sparalog.TraceLevel, "", false, args...)
}

// Printf sends to Default logger info stream using the same fmt.Printf() interface.
func Printf(format string, args ...interface{}) {
	Default.Logf(sparalog.InfoLevel, "", false, format, args...)
}

// Print sends to Default logger info stream using the same fmt.Print() interface.
func Print(args ...interface{}) {
	Default.Log(sparalog.InfoLevel, "", false, args...)
}

// ResetWriters reset the writers for all the levels to an optional default writer.
func ResetWriters(defaultW sparalog.Writer) {
	Default.ResetWriters(defaultW)
}

// ResetLevelWriters remove all level's writers and reset to an optional default writer.
func ResetLevelWriters(level sparalog.Level, defaultW sparalog.Writer) {
	Default.ResetLevelWriters(level, defaultW)
}

// ResetLevelsWriters remove specific levels writers and reset to an optional default writer.
func ResetLevelsWriters(levels []sparalog.Level, defaultW sparalog.Writer) {
	Default.ResetLevelsWriters(levels, defaultW)
}

// AddWriter add a writer to all levels.
// id is optional, but useful for RemoveWriter().
func AddWriter(w sparalog.Writer, id sparalog.WriterID) {
	Default.AddWriter(w, id)
}

// AddLevelWriter add a writer to a level.
// id is optional, but useful for RemoveWriter().
func AddLevelWriter(level sparalog.Level, w sparalog.Writer, id sparalog.WriterID) {
	Default.AddLevelWriter(level, w, id)
}

// AddLevelsWriter add a writer to several levels.
// id is optional, but useful for RemoveWriter().
func AddLevelsWriter(levels []sparalog.Level, w sparalog.Writer, id sparalog.WriterID) {
	Default.AddLevelsWriter(levels, w, id)
}

// RemoveWriter delete a specific writer from level.
func RemoveWriter(level sparalog.Level, id sparalog.WriterID) {
	Default.RemoveWriter(level, id)
}

// Mute mute/unmute a specific level.
func Mute(level sparalog.Level, state bool) {
	Default.Mute(level, state)
}

// EnableStacktrace enable stacktrace for a specific level.
func EnableStacktrace(level sparalog.Level, state bool) {
	Default.EnableStacktrace(level, state)
}

// How many top calls to skip from the stack trace.
var skipStackCalls = 5
