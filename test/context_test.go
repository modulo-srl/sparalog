package test

// TODO

import (
	"errors"
	"fmt"
	"testing"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
)

// TODO
func TestContextData(t *testing.T) {
}

// TODO
func TestContextTags(t *testing.T) {
	defer logs.Done()
	logs.Open()

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

// TODO
func TestPrefix(t *testing.T) {
}

// TODO
func TestContextPrefix(t *testing.T) {
	defer logs.Done()
	logs.Open()

	logs.SetContextPrefix("%s", []string{"module"})

	sublog := logs.NewAliasLogger()
	sublog.SetContextData("module", "child")

	i := sublog.NewItem(sparalog.InfoLevel, "test")
	s := i.RenderPrefix()
	//s := i.ToString(false, false)
	fmt.Println(s)

	//os.Exit(1)
}
