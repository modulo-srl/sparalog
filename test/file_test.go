package test

import (
	"testing"

	"github.com/modulo-srl/sparalog/logs"
)

func TestFile(t *testing.T) {
	defer logs.Done()

	logs.Init("sparalog-test")

	w, err := logs.NewFileWriter("test.log")
	if err != nil {
		t.Fatal(err)
	}
	logs.ResetAllWriters(w)

	logs.StartPanicWatcher()

	logs.Error("test")

	// TODO load file and compare
}
