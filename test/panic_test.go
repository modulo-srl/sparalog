package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
	"github.com/modulo-srl/sparalog/writers"
)

func TestMainPanic(t *testing.T) {
	return

	sparalog.InitUnitTest()

	sparalog.Start()
	defer sparalog.Stop()

	w := writers.NewCallbackWriter(
		func(item *logs.Item) error {
			if item.Level == logs.FatalLevel {
				os.Exit(0)
			}
			return nil
		},
	)
	//logs.ResetWriters(w)
	logs.AddWriter(w)

	makePanic()
	t.Fatal("panic not logged")
}

/*func TestGoroutinePanic(t *testing.T) {
	logs.Init("sparalog-test")
	defer logs.Done()

	w := logs.NewCallbackWriter(
		func(item *logs.Item) error {
			if item.Level() == logs.FatalLevel {
				// Received by the parent process,
				// that should exits shortly (see logger.Fatal()),
				// so change the exit code to success.
				logs.FatalExitCode = 0
			}
			return nil
		},
	)
	logs.ResetWriters(w)

	logs.StartPanicWatcher()

	go func() {
		makePanic()
	}()
	time.Sleep(time.Second) // assure goroutine execution
}*/

func makePanic() {
	i := 0
	i = 1 / i

	fmt.Print(i)
}
