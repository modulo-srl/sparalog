package logs

// Funzioni interne di gestione logger di default.

// Logger globale di default.
var defaultLogger *Logger

// Dispatcher globale.
var globalDispatcher *dispatcher

// Per quali livelli lo stacktrace è abilitato.
var levelsStackTrace [LevelsCount]bool

// Invocata da sparalog.init()
// Inizializza la libreria allocando il logger di default
// e associando un writer di default (di tipo Stdout).
func InitDefaultLogger(defaultWriter Writer) {
	defaultLogger = &Logger{
		// Skippa una chiamata in più dato che il logger di default
		// viene invocato da funzioni wrapper.
		stackCallsToSkip: 1,
	}

	// Inizializza il dispatcher.
	globalDispatcher = newDispatcher(defaultWriter)

	// Abilita lo stacktrace per i soli livelli fatal, error.
	EnableLevelsStackTrace([]Level{FatalLevel, ErrorLevel})

	globalDispatcher.Mute(DebugLevel, true)
}

// Invocata da sparalog.Start()
// Avvia il logger di default.
func StartDefaultLogger() {
	globalDispatcher.Start()
}

// Invocata da sparalog.Stop()
// Gestisce i panic (della routine corrente),
// termina il dispatcher attendendo gentilmente il termine dei writer asincroni.
// Va chiamata al termine del main.
func StopDefaultLogger() {
	err := recover()
	if err != nil {
		Fatal(err)
	}

	if globalDispatcher == nil {
		return
	}

	globalDispatcher.Stop()
	globalDispatcher = nil
}
