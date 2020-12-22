package test

import (
	"testing"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
)

func TestRootLoggerMethods(t *testing.T) {
	logs.Init("sparalog-test")

	var traced string

	w := logs.NewCallbackWriter(
		func(item sparalog.Item) error {
			traced = item.ToString(false, true)
			return nil
		},
	)
	logs.ResetWriters(w)

	logs.Error("asd")
	traced1 := traced
	logs.Default.Error("asd")
	traced2 := traced
	if traced1 != traced2 {
		t.Fatal("Error() not same")
	}

	logs.Warn("asd")
	traced1 = traced
	logs.Default.Warn("asd")
	traced2 = traced
	if traced1 != traced2 {
		t.Fatal("Warn() not same")
	}

	logs.Info("asd")
	traced1 = traced
	logs.Default.Info("asd")
	traced2 = traced
	if traced1 != traced2 {
		t.Fatal("Info() not same")
	}

	logs.Debug("asd")
	traced1 = traced
	logs.Default.Debug("asd")
	traced2 = traced
	if traced1 != traced2 {
		t.Fatal("Debug() not same")
	}

	logs.Trace("asd")
	traced1 = traced
	logs.Default.Trace("asd")
	traced2 = traced
	if traced1 != traced2 {
		t.Fatal("Trace() not same")
	}

	logs.Print("asd")
	traced1 = traced
	logs.Default.Print("asd")
	traced2 = traced
	if traced1 != traced2 {
		t.Fatal("Print() not same")
	}
}
