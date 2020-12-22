package test

import (
	"errors"
	"strings"
	"testing"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/item"
	"github.com/modulo-srl/sparalog/logs"
)

func TestStacktraceLogger(t *testing.T) {
	defer logs.Done()
	logs.Init("sparalog-test")

	ws := logs.NewCallbackWriter(
		func(item sparalog.Item) error {
			checkStacktrace(item, t)
			return nil
		},
	)
	logs.ResetLevelWriters(sparalog.ErrorLevel, ws)

	logs.Error("test-lib")

	logs.Default.Error("test-logger")

	i := logs.NewError(errors.New("test-lib-item"))
	checkStacktrace(i, t)

	i = logs.Default.NewError(errors.New("test-logger-item"))
	checkStacktrace(i, t)

	i = item.NewError(0, errors.New("test-logitem"))
	checkStacktrace(i, t)
}

// check if the stacktrace starts with current testing function.
func checkStacktrace(item sparalog.Item, t *testing.T) {
	lines := strings.Split(item.StackTrace(), "\n")

	if len(lines) < 2 {
		t.Errorf("No stacktrace for \"%s\"\n%s", item.Line(), item.StackTrace())
		return
	}

	i := strings.Index(lines[1], "test.Test")

	if i < 0 {
		t.Errorf("Invalid stacktrace for \"%s\"\n%s", item.Line(), item.StackTrace())
	}
}
