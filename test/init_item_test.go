package test

import (
	"testing"
	"time"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
	"github.com/modulo-srl/sparalog/writers"
)

func TestInitItem(t *testing.T) {
	sparalog.InitUnitTest()

	cw := writers.NewCallbackWriter(testInitItemCallback)
	logs.AddWriter(cw)

	sparalog.Start()
	defer sparalog.Stop()

	logs.SetInitItemFunc(testInitItem)

	logs.Info("test")

	time.Sleep(50 * time.Millisecond)

	if lastPrefix != "test-prefix" {
		t.Fatal("prefix error: " + lastPrefix)
	}

}

func testInitItem(item *logs.Item) {
	item.Prefix += "test-prefix"
}

var lastPrefix string

func testInitItemCallback(i *logs.Item) error {
	lastPrefix = i.Prefix
	return nil
}
