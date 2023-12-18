package main

import (
	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
)

func main() {
	// Per default va tutto su stdout/stderror
	sparalog.Start()
	defer sparalog.Stop()

	logs.Infof("%s", "info!")

	logs.Error("error!")
}
