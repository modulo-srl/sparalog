package test

import (
	"errors"
	"strings"
	"testing"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
	"github.com/modulo-srl/sparalog/writers"
)

func TestStacktraceLogger(t *testing.T) {
	sparalog.InitUnitTest()

	sparalog.Start()
	defer sparalog.Stop()

	ws := writers.NewCallbackWriter(
		func(item *logs.Item) error {
			checkStacktrace(item, t)
			return nil
		},
	)
	logs.ResetLevelWriters(logs.ErrorLevel, ws)

	logs.Error("test-lib")

	i := logs.NewItem(logs.ErrorLevel, errors.New("test-lib-item"))
	checkStacktrace(i, t)

	i = logs.NewItemf(logs.ErrorLevel, "%s", "test-lib-itemf")
	checkStacktrace(i, t)

	i = logs.NewErrorItem(errors.New("test-lib-error-item"))
	checkStacktrace(i, t)

	i = logs.NewErrorItemf("%s", "test-lib-errorf-item")
	checkStacktrace(i, t)

	logger := logs.NewLogger("")

	logger.Error("test-logger")

	i = logger.NewItem(logs.ErrorLevel, errors.New("test-logger-item"))
	checkStacktrace(i, t)

	i = logger.NewItemf(logs.ErrorLevel, "%s", "test-logger-itemf")
	checkStacktrace(i, t)

	i = logger.NewErrorItem(errors.New("test-logger-error-item"))
	checkStacktrace(i, t)

	i = logger.NewErrorItemf("%s", "test-logger-errorf-item")
	checkStacktrace(i, t)
}

// check if the stacktrace starts with current testing function.
func checkStacktrace(item *logs.Item, t *testing.T) {
	lines := strings.Split(item.StackTrace, "\n")

	if len(lines) < 2 {
		t.Errorf("No stacktrace for \"%s\"\n%s", item.Message, item.StackTrace)
		return
	}

	i := strings.Index(lines[1], "test.TestStacktraceLogger")

	if i < 0 {
		t.Errorf("Invalid stacktrace for \"%s\"\n%s", item.Message, item.StackTrace)
	}
}
