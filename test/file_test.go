package test

import (
	"testing"

	"github.com/modulo-srl/sparalog/logs"
)

func TestFile(t *testing.T) {
	logs.Init("sparalog-test")
	defer logs.Done()

	w, err := logs.NewFileWriter("test.log")
	if err != nil {
		t.Fatal(err)
	}
	logs.ResetWriters(w)

	logs.StartPanicWatcher()

	logs.Error("test")

	// TODO load file and compare
}
