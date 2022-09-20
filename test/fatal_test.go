package test

import (
	"testing"

	"github.com/modulo-srl/sparalog/logs"
)

// TODO
func TestFatal(t *testing.T) {
	logs.Open()
	//logs.StartPanicWatcher()

	//logs.Fatal("test fatal")
}

// TODO
func TestError(t *testing.T) {
	logs.Open()

	logs.Error("test error")

	//t.Error("")
}
