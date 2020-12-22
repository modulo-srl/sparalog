package test

import (
	"testing"

	"github.com/modulo-srl/sparalog/logs"
)

// TODO
func TestFatal(t *testing.T) {
	logs.Init("sparalog-test")
	//logs.StartPanicWatcher()

	//logs.Fatal("test fatal")
}

// TODO
func TestError(t *testing.T) {
	logs.Init("sparalog-test")

	logs.Error("test error")

	//t.Error("")
}
