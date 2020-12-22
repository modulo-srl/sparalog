package logs

// Logging domain controller.
// See root package for the model and interfaces.

import (
	"os"
	"strings"

	"github.com/mitchellh/panicwrap"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/env"
	"github.com/modulo-srl/sparalog/logger"
	"github.com/modulo-srl/sparalog/writers"
)

// Default is the default global logger.
var Default sparalog.Logger

// DefaultDispatcher is the default global dispatcher.
var DefaultDispatcher sparalog.Dispatcher

// DefaultStdoutWriter is the default global writer to stdout.
var DefaultStdoutWriter sparalog.Writer

// Init initialize the library - it should be called at main start.
// programName = program name and version.
func Init(programName string) {
	if Default != nil {
		return
	}
	closed = false

	env.Init(programName)

	loggers = make([]sparalog.Logger, 0, 1)

	DefaultStdoutWriter = writers.NewStdoutWriter()
	Default = NewLogger()
	DefaultDispatcher = Default.(*logger.Logger).Dispatcher
}

// StartPanicWatcher starts a supervisor that monitors panics in all goroutines.
// Since the supervision is made starting a parent + child processes:
// - Call the function after all writers initialization, or at least after fatal level initialization;
// - This function should not to be called in debugging sessions.
func StartPanicWatcher() {
	exitStatus, err := panicwrap.Wrap(&panicwrap.WrapConfig{
		Handler:   panicHandler,
		HidePanic: true,
	})
	if err != nil {
		panic(err)
	}

	// If exitStatus >= 0, then we're the parent process and the panicwrap
	// re-executed ourselves and completed. Just exit with the proper status.
	if exitStatus >= 0 {
		os.Exit(exitStatus)
	}

	// Otherwise, exitStatus < 0 means we're the child. Continue executing as
	// normal...
}

// Done manages (current routine) panics and closes all the loggers,
// waiting gently for the async writers.
// It should be called at main exit.
func Done() {
	err := recover()
	if err != nil {
		Fatal(err)
	}

	if closed {
		return
	}
	closed = true

	for _, l := range loggers {
		l.Close()
	}

	Default = nil
	DefaultDispatcher = nil
	DefaultStdoutWriter = nil
}

var loggers []sparalog.Logger
var closed bool

// output contains the full output (including stack traces) of the child panic.
func panicHandler(output string) {
	var st string

	// Strip the stack trace.
	i := strings.Index(output, "\n\n")
	if i >= 0 {
		st = "STACKTRACE (by watcher): " + strings.TrimSpace(output[i+2:])
		output = output[:i]
	}

	item := Default.NewItem(sparalog.FatalLevel, output)
	item.SetStackTrace(st)
	item.Log()
}
