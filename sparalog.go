package sparalog

// Funzioni global di start e stop.
// Per il resto fare riferimento ai package `logs` e `writers`.

import (
	"github.com/modulo-srl/sparalog/logs"
	"github.com/modulo-srl/sparalog/writers"
)

// Avvia il logger di default.
// Va chiamata una volta che i writer sono stati inizializzati e associati.
func Start() {
	logs.StartDefaultLogger()
}

// Termina il sistema di logging.
// Gestisce i panic (della routine corrente), termina i logger
// attendendo gentilmente il termine dei writer asincroni.
// Va chiamata al termine del main.
func Stop() {
	logs.StopDefaultLogger()
}

// Inizializza il sistema di logging.
// Alloca un writer di tipo stdout e lo passa alla funzione
// di inizializzazione del logger di default.
func init() {
	initSparalog()
}

// Da chiamare prima di ogni unit test.
func InitUnitTest() {
	initSparalog()
}

func initSparalog() {
	w := writers.NewStdoutWriter()
	logs.InitDefaultLogger(w)
}
