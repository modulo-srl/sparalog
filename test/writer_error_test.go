package test

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
)

func TestWriterError(t *testing.T) {
	var forwardLoggeds int
	var mu sync.Mutex

	// Default writer.
	ws := logs.NewCallbackWriter(
		func(item sparalog.Item) error {
			fmt.Println(item.ToString(true, true))

			mu.Lock()
			defer mu.Unlock()
			forwardLoggeds++
			return nil
		},
	)
	logs.ResetLevelWriters(sparalog.ErrorLevel, ws)

	wa := logs.NewCallbackWriter(
		func(item sparalog.Item) error {
			return errors.New("feedback error")
		},
	)
	logs.AddWriter(wa)

	logs.Open()
	defer logs.Done()

	logs.Info("test writer error")

	time.Sleep(time.Second) // assure writers forward channel processing.

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
	var forwardLoggeds int
	var mu sync.Mutex

	// Default writer.
	ws := logs.NewCallbackAsyncWriter(
		func(item sparalog.Item) error {
			fmt.Println("*** CATCHED! ***")
			fmt.Println(item.ToString(true, true))

			mu.Lock()
			defer mu.Unlock()
			forwardLoggeds++
			return nil
		},
	)
	logs.ResetLevelWriters(sparalog.ErrorLevel, ws)

	wa := logs.NewCallbackAsyncWriter(
		func(item sparalog.Item) error {
			return errors.New("feedback error")
		},
	)
	logs.AddWriter(wa)

	logs.Open()
	defer logs.Done()

	logs.Info("test async writer error")

	time.Sleep(time.Second) // assure writers forward channel processing.

	mu.Lock()
	defer mu.Unlock()

	switch {
	case forwardLoggeds == 0:
		t.Error("forward not logged")
	case forwardLoggeds > 1:
		t.Error("too many forwards. deadlock?")
	}
}
