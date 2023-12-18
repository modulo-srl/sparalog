package test

import (
	"testing"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
)

// TODO
func TestFatal(t *testing.T) {
	sparalog.InitUnitTest()

	sparalog.Start()
	//logs.StartPanicWatcher()

	//logs.Fatal("test fatal")
}

// TODO
func TestError(t *testing.T) {
	sparalog.InitUnitTest()

	sparalog.Start()

	logs.Error("test error")

	//t.Error("")
}
