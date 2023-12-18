package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
	"github.com/modulo-srl/sparalog/writers"
)

func TestPayloadData(t *testing.T) {
	sparalog.InitUnitTest()

	cw := writers.NewCallbackWriter(myCallback)
	logs.AddWriter(cw)

	sparalog.Start()
	defer sparalog.Stop()

	logs.SetPayload("by main", "")

	i := logs.NewItem(logs.InfoLevel, "test main")
	i.SetPayload("by main item", "")
	logs.LogItem(i)

	time.Sleep(20 * time.Millisecond)
	if err := checkPayload([]string{"by main item", "by main item"}); err != nil {
		t.Fatal(err)
	}

	m := module{}
	m.init()

	time.Sleep(20 * time.Millisecond)
	if err := checkPayload([]string{"by main", "by module"}); err != nil {
		t.Fatal(err)
	}

	m.test()

	time.Sleep(20 * time.Millisecond)
	if err := checkPayload([]string{"by main", "by module", "by module item"}); err != nil {
		t.Fatal(err)
	}
}

// ultimo payload ricevuto
var lastPayload map[string]any

func checkPayload(keys []string) error {
	if lastPayload == nil {
		return fmt.Errorf("expected: %v, got: %v", keys, nil)
	}

	for _, key := range keys {
		_, ok := lastPayload[key]
		if !ok {
			return fmt.Errorf("expected: %v, got: %v", keys, lastPayload)
		}
	}

	return nil
}

type module struct {
	log *logs.Logger
}

func (m *module) init() {
	m.log = logs.NewLogger("my module")
	m.log.SetPayload("by module", "")

	m.log.Info("initialized")
}

func (m *module) test() {
	i := m.log.NewItem(logs.InfoLevel, "test")
	i.SetPayload("by module item", "")
	m.log.LogItem(i)
}

func myCallback(i *logs.Item) error {
	lastPayload = i.Payload
	return nil
}
