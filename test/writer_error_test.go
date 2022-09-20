package test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
)

func TestWriterError(t *testing.T) {
	logs.Open()
	defer logs.Done()

	var forwardLoggeds int

	// Default writer.
	ws := logs.NewCallbackWriter(
		func(item sparalog.Item) error {
			fmt.Println(item.ToString(true, true))
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
	logs.AddWriter(wa, "")

	logs.Info("test writer error")

	time.Sleep(time.Second) // assure writers forward channel processing.

	switch {
	case forwardLoggeds == 0:
		t.Error("forward not logged")
	case forwardLoggeds > 1:
		t.Error("too many forwards. deadlock?")
	}
}

func TestAsyncWriterError(t *testing.T) {
	logs.Open()
	defer logs.Done()

	var forwardLoggeds int

	// Default writer.
	ws := logs.NewCallbackAsyncWriter(
		func(item sparalog.Item) error {
			fmt.Println("*** CATCHED! ***")
			fmt.Println(item.ToString(true, true))
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
	logs.AddWriter(wa, "")

	logs.Info("test async writer error")

	time.Sleep(time.Second) // assure writers forward channel processing.

	switch {
	case forwardLoggeds == 0:
		t.Error("forward not logged")
	case forwardLoggeds > 1:
		t.Error("too many forwards. deadlock?")
	}
}
