package main

import (
	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/logs"
	"github.com/modulo-srl/sparalog/writers"
)

func main() {
	// Writer su file.
	wf, err := writers.NewFileWriter("/tmp/example.log")
	if err != nil {
		panic(err)
	}

	// Resetta i writers (no stdout) impostando wf come quello di default per tutti i livelli.
	logs.ResetWriters(wf)

	// Logga il livello di debug verso stdout.
	wstd := writers.NewStdoutWriter()
	logs.AddLevelWriter(logs.DebugLevel, wstd)

	// Logga i livelli critici anche verso Syslog.
	wsys := writers.NewSyslogWriter("example")
	logs.AddLevelsWriter(logs.CriticalLevels, wsys)

	// Logga i fatali anche verso Telegram.
	wt := writers.NewTelegramWriter("apikey", 1234567890)
	logs.AddLevelWriter(logs.FatalLevel, wt)

	// TCP alla porta 6006, che si auto abilita quando ci si connette.
	wtcp, _ := writers.NewTCPWriter("", 6006, false, func(state bool) {
		// true = è entrato il primo client; false = è uscito l'ultimo client
		if state {
			logs.Info("enable tracing log to TCP")
		} else {
			logs.Info("disable tracing log to TCP")
		}
	})
	logs.AddWriter(wtcp)

	sparalog.Start()
	defer sparalog.Stop()

	// ...

	// Abilita lo stacktrace per i livelli warning, error, fatal.
	logs.EnableLevelsStackTrace(logs.CriticalLevels)

	// Muta temporaneamente il livello Debug.
	logs.Mute(logs.DebugLevel, true)

	// ...
}
