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

	return logger.NewChild(nil, tag)
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

// ResetAllWriters reset the writers for all the levels to an optional default writer.
func ResetAllWriters(defaultW sparalog.Writer) {
	Default.ResetAllWriters(defaultW)
}

// ResetWriters remove all level's writers and reset to an optional default writer.
func ResetWriters(level sparalog.Level, defaultW sparalog.Writer) {
	Default.ResetWriters(level, defaultW)
}

// AddWriter add a writer to a level.
// id is optional, but useful for RemoveWriter().
func AddWriter(level sparalog.Level, w sparalog.Writer, id sparalog.WriterID) {
	Default.AddWriter(level, w, id)
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
