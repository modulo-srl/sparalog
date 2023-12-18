package test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
	"github.com/modulo-srl/sparalog/writers"
)

func TestFile(t *testing.T) {
	sparalog.InitUnitTest()

	sparalog.Start()
	defer sparalog.Stop()

	os.Remove("test.log")

	w, err := writers.NewFileWriter("test.log")
	if err != nil {
		t.Fatal(err)
	}
	logs.ResetWriters(w)

	logs.Error("test-file")

	time.Sleep(50 * time.Millisecond)

	bb, err := os.ReadFile("test.log")
	if err != nil {
		t.Fatal(err)
	}

	s := string(bb)
	if strings.Index(s, "test-file") <= 0 {
		t.Fatal("mismatch: ", s)
	}
}
