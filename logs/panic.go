package logs

// Gestione del watcher dei panic.

import (
	"os"
	"strings"
	"time"

	"github.com/mitchellh/panicwrap"
)

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

// Handler per i panic.
//   - output: contiene l'intero output (compreso di stacktrace)
//     del panic del processo figlio.
func panicHandler(output string) {
	var st string

	time.Sleep(time.Second * 3)

	// Aggiorna lo stacktrace.
	i := strings.Index(output, "\n\n")
	if i >= 0 {
		st = "STACKTRACE (by panic watcher): " + strings.TrimSpace(output[i+2:])
		output = output[:i]
	}

	item := NewItem(FatalLevel, output)
	item.StackTrace = st
	defaultLogger.LogItem(item)

	globalDispatcher.Stop()
}
