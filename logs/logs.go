package logs

// Funzioni e variabili generali.

// Interfaccia del writer usata dal logger.
type Writer interface {
	ID() string

	Write(*Item)

	Start() error
	Stop()
	SetFeedbackChan(chan *Item)
}

// Exit Code generato da Fatal() e Fatalf().
var FatalExitCode = 1

// Funzione di inizializzazione item appena dopo la sua allocazione,
// permette di settare ulteriormente in modo custom specifiche proprietà.
type InitItemF func(*Item)

// Alloca un nuovo logger dotato di prefisso e di differente payload,
// utilizzabile come componente in una struct o in un package.
// Eredita una copia del payload dal logger di default.
// - prefix: prefisso di default che comparirà nelle relative loggate.
func NewLogger(prefix string) *Logger {
	return newAliasLogger(defaultLogger, prefix)
}

// Setta un valore del payload di default del logger di default.
func SetPayload(name string, value any) {
	defaultLogger.SetPayload(name, value)
}

// Imposta una funzione di inizializzazione per ogni item allocato dal logger di default.
func SetInitItemFunc(f InitItemF) {
	defaultLogger.initItemF = f
}
