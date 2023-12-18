package main

import (
	"fmt"
	"os"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
	"github.com/modulo-srl/sparalog/writers"
)

func main() {
	// Logga i fatali anche su file.
	os.Remove("test.log")
	w, err := writers.NewFileWriter("test.log")
	if err != nil {
		panic(err)
	}

	logs.AddLevelWriter(logs.FatalLevel, w)

	w2 := writers.NewStdoutWriter()
	logs.AddLevelWriter(logs.FatalLevel, w2)

	wt := writers.NewTelegramWriter("429865516:AAGNWELnsfY5N-la2lsb2FC2WyHvNo8dqrI", -495679821)
	logs.AddLevelWriter(logs.FatalLevel, wt)

	sparalog.Start()
	//defer sparalog.Stop()

	// Avvia il watcher.
	// note:
	// - i writer, o comunque quelli relativi al livello fatal, devono essere settati prima di questa chiamata.
	// - l'avvio del watcher va evitato per le sessioni di debug
	logs.StartPanicWatcher()

	go func() {
		i := 0
		i = 1 / i // andrà in panic e verrà loggato

		fmt.Print(i)
	}()
}
