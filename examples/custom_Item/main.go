package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
	"github.com/modulo-srl/sparalog/writers"
)

func myInitItem(item *logs.Item) {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	goID := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]

	p := "#" + goID
	if item.Prefix != "" {
		item.Prefix = p + " " + item.Prefix
	} else {
		item.Prefix = p
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	item.SetPayload("ram", m.TotalAlloc)
	item.SetPayload("numGC", m.NumGC)
}

func myCallback(i *logs.Item) error {
	fmt.Println("payload", i.Payload)
	return nil
}

func main() {
	cw := writers.NewCallbackWriter(myCallback)
	logs.AddWriter(cw)

	logs.SetInitItemFunc(myInitItem)

	sparalog.Start()
	defer sparalog.Stop()

	logs.SetPayload("foo", "bar")

	logs.Info("test")
}
