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
	logs.Init("sparalog-test")
	defer logs.Done()

	var forwardLoggeds int

	// Default writer.
	ws := logs.NewCallbackWriter(
		func(item sparalog.Item) error {
			fmt.Println(item.String(true, true))
			forwardLoggeds++
			return nil
		},
	)
	logs.ResetLevelWriters(sparalog.ErrorLevel, ws)
	logs.EnableStacktrace(sparalog.ErrorLevel, true)

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
	logs.Init("sparalog-test")
	defer logs.Done()

	var forwardLoggeds int

	// Default writer.
	ws := logs.NewCallbackAsyncWriter(
		func(item sparalog.Item) error {
			fmt.Println(item.String(true, true))
			forwardLoggeds++
			return nil
		},
	)
	logs.ResetLevelWriters(sparalog.ErrorLevel, ws)
	logs.EnableStacktrace(sparalog.ErrorLevel, true)

	wa := logs.NewCallbackAsyncWriter(
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
