package test

// TODO

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
)

func TestContextData(t *testing.T) {
}

func TestContextTags(t *testing.T) {
	defer logs.Done()
	logs.Init("sparalog-test")

	logs.SetContextTag("peerID", "pid")
	logs.SetContextData("foo", "bar")
	logs.SetContextPrefix("%s: %s", []string{"mh", "peerID"})

	logs.Error("this is an error: ", errors.New("error"))

	i := logs.NewErrorf("this %s is an error", errors.New("error"))
	i.SetPrefix("%s = %s", []string{"foo", "peerID"})
	i.SetTag("foo", "bar")
	i.SetTag("peerID", "pid2")
	i.Log()

	//os.Exit(1)
}

func TestPrefix(t *testing.T) {
}

func TestFingerprint(t *testing.T) {
	defer logs.Done()
	logs.Init("sparalog-test")

	/*ws := logs.NewCallbackWriter(
		func(item sparalog.Item) error {
			fmt.Println(item.ToString(true, true))
			forwardLoggeds++
			return nil
		},
	)
	logs.ResetLevelWriters(sparalog.ErrorLevel, ws)*/

	i := logs.NewError(errors.New("test-error"))
	fp := i.Fingerprint()

	if len(fp) < 32 {
		t.Fatal("Fingerprint error:", fp)
	}
}

func TestContextPrefix(t *testing.T) {
	defer logs.Done()
	logs.Init("sparalog-test")

	logs.SetContextPrefix("%s", []string{"module"})

	sublog := logs.NewAliasLogger()
	sublog.SetContextData("module", "child")

	i := sublog.NewItem(sparalog.InfoLevel, "test")
	s := i.RenderPrefix()
	//s := i.ToString(false, false)
	fmt.Println(s)

	os.Exit(1)
}
