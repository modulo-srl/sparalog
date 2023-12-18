package test

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
	"github.com/modulo-srl/sparalog/writers"
)

func TestWriterError(t *testing.T) {
	sparalog.InitUnitTest()

	var forwardLoggeds int
	var mu sync.Mutex

	// Default writer.
	ws := writers.NewCallbackWriter(
		func(item *logs.Item) error {
			fmt.Println(item.ToString(true, true))

			mu.Lock()
			defer mu.Unlock()
			forwardLoggeds++
			return nil
		},
	)
	logs.ResetLevelWriters(logs.ErrorLevel, ws)

	wa := writers.NewCallbackWriter(
		func(item *logs.Item) error {
			return errors.New("feedback error")
		},
	)
	logs.AddWriter(wa)

	sparalog.Start()
	defer sparalog.Stop()

	logs.Info("test writer error")

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	switch {
	case forwardLoggeds == 0:
		t.Error("forward not logged")
	case forwardLoggeds > 1:
		t.Error("too many forwards. deadlock?")
	}
}

func TestAsyncWriterError(t *testing.T) {
	sparalog.InitUnitTest()

	var forwardLoggeds int
	var mu sync.Mutex

	// Default writer.
	ws := writers.NewCallbackAsyncWriter(
		func(item *logs.Item) error {
			fmt.Println("*** CATCHED! ***")
			fmt.Println(item.ToString(true, true))

			mu.Lock()
			defer mu.Unlock()
			forwardLoggeds++
			return nil
		},
	)
	logs.ResetLevelWriters(logs.ErrorLevel, ws)

	wa := writers.NewCallbackAsyncWriter(
		func(item *logs.Item) error {
			return errors.New("feedback error")
		},
	)
	logs.AddWriter(wa)

	sparalog.Start()
	defer sparalog.Stop()

	logs.Info("test async writer error")

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	switch {
	case forwardLoggeds == 0:
		t.Error("forward not logged")
	case forwardLoggeds > 1:
		t.Error("too many forwards. deadlock?")
	}
}
