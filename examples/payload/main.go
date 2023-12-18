package main

import (
	"fmt"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
	"github.com/modulo-srl/sparalog/writers"
)

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
	i.SetPayload("by item module", "")
	m.log.LogItem(i)
}

func myCallback(i *logs.Item) error {
	fmt.Println("payload", i.Payload)
	return nil
}

func main() {
	cw := writers.NewCallbackWriter(myCallback)
	logs.AddWriter(cw)

	sparalog.Start()
	defer sparalog.Stop()

	logs.SetPayload("by main", "")

	i := logs.NewItem(logs.InfoLevel, "test main")
	i.SetPayload("by main item", "")
	logs.LogItem(i)

	m := module{}
	m.init()

	m.test()
}
