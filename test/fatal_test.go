package test

import (
	"testing"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
)

func TestFatal(t *testing.T) {
	logs.Init("sparalog-test")
	logs.StartPanicWatcher()

	logs.Fatal("test fatal")
}

func TestError(t *testing.T) {
	logs.Init("sparalog-test")

	logs.EnableStacktrace(sparalog.ErrorLevel, true)
	logs.Error("test error")

	t.Error("")
}
