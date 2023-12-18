package main

import (
	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
)

type module struct {
	log *logs.Logger
}

func (m *module) init() {
	m.log = logs.NewLogger("my module")

	m.log.Info("initialized")
}

func main() {
	sparalog.Start()
	defer sparalog.Stop()

	m := module{}
	m.init()
}
