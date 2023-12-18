package logs

// Logger.

import (
	"fmt"
	"sync"
)

type Logger struct {
	initItemF InitItemF

	prefix string

	// Quante chiamate dello stackTrace escludere di default.
	stackCallsToSkip int

	muPayload sync.RWMutex
	payload   map[string]any
}

// Alloca un nuovo logger.
func newAliasLogger(logger *Logger, prefix string) *Logger {
	l := Logger{
		prefix:           prefix,
		payload:          logger.getPayloadCopy(), // riceve una copia del payload se il logger padre ne è provvisto
		initItemF:        logger.initItemF,
		stackCallsToSkip: 0,
	}

	return &l
}

// Logga a livello fatale effetuando anche os.Exit(FatalExitCode)
func (l *Logger) Fatal(args ...any) {
	l.log(FatalLevel, args...)
}

// Logga a livello fatale effetuando anche os.Exit(FatalExitCode)
func (l *Logger) Fatalf(format string, args ...any) {
	l.log(FatalLevel, fmt.Sprintf(format, args...))
}

// Logga a livello errore.
func (l *Logger) Error(args ...any) {
	l.log(ErrorLevel, args...)
}

// Logga a livello errore.
func (l *Logger) Errorf(format string, args ...any) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...))
}

// Logga a livello warning.
func (l *Logger) Warning(args ...any) {
	l.log(WarningLevel, args...)
}

// Logga a livello warning.
func (l *Logger) Warningf(format string, args ...any) {
	l.log(WarningLevel, fmt.Sprintf(format, args...))
}

// Logga a livello info.
func (l *Logger) Info(args ...any) {
	l.log(InfoLevel, args...)
}

// Logga a livello info.
func (l *Logger) Infof(format string, args ...any) {
	l.log(InfoLevel, fmt.Sprintf(format, args...))
}

// Logga a livello debug.
func (l *Logger) Debug(args ...any) {
	l.log(DebugLevel, args...)
}

// Logga a livello debug.
func (l *Logger) Debugf(format string, args ...any) {
	l.log(DebugLevel, fmt.Sprintf(format, args...))
}

// Genera un nuovo item di livello specifico.
// Eredita una copia del payload dal logger che può essere ulteriormente customizzata.
func (l *Logger) NewItem(level Level, args ...any) *Item {
	return l.newSelfItem(level, 1, fmt.Sprint(args...))
}

// Genera un nuovo item di livello specifico.
// Eredita una copia del payload dal logger che può essere ulteriormente customizzata.
func (l *Logger) NewItemf(level Level, format string, args ...any) *Item {
	return l.newSelfItem(level, 1, fmt.Sprintf(format, args...))
}

// Genera un nuovo item di livello errore.
// Eredita una copia del payload dal logger che può essere ulteriormente customizzata.
func (l *Logger) NewErrorItem(err error) *Item {
	return l.newSelfItem(ErrorLevel, 1, err.Error())
}

// Genera un nuovo item di livello errore.
// Eredita una copia del payload dal logger che può essere ulteriormente customizzata.
func (l *Logger) NewErrorItemf(format string, a ...any) *Item {
	return l.newSelfItem(ErrorLevel, 1, fmt.Errorf(format, a...).Error())
}

// Setta un valore del payload di default.
func (l *Logger) SetPayload(key string, value any) {
	l.muPayload.Lock()
	defer l.muPayload.Unlock()

	if l.payload == nil {
		l.payload = make(map[string]any)
	}

	l.payload[key] = value
}

// Logga un item precedentemente generato.
func (l *Logger) LogItem(item *Item) {
	if !globalDispatcher.CanDispatch(item.Level) {
		return
	}

	globalDispatcher.Dispatch(item)
}

// Imposta una funzione di inizializzazione per ogni item allocato dal logger.
func (l *Logger) SetInitItemFunc(f InitItemF) {
	l.initItemF = f
}

// Ritorna una copia del payload di default.
func (l *Logger) getPayloadCopy() map[string]any {
	l.muPayload.RLock()
	defer l.muPayload.RUnlock()

	if l.payload == nil {
		return nil
	}

	payload := make(map[string]any)
	for k, v := range l.payload {
		payload[k] = v
	}

	return payload
}

// Genera un nuovo item cedibile allo strato applicativo,
// utilizzata dai vari logger.NewIem*() e logger.NewError*().
// Eredita una copia del payload dal logger che può essere ulteriormente customizzata
// ed eventualmente inizializza l'item con la funzione custom.
func (l *Logger) newSelfItem(level Level, stackCallsToSkip int, msg string) *Item {
	item := newItem(level, l.prefix, msg, l.stackCallsToSkip+stackCallsToSkip+1)

	// Assegna una copia del payload del logger se ne è provvisto.
	item.Payload = l.getPayloadCopy()

	if l.initItemF != nil {
		l.initItemF(item)
	}

	return item
}

// Logga in uno specifico livello - entry point per tutti gli helper che loggano; thread safe.
// Non fa nulla se il livello è mutato.
func (l *Logger) log(level Level, args ...any) {
	if !globalDispatcher.CanDispatch(level) {
		return
	}

	item := newItem(level, l.prefix, fmt.Sprint(args...), l.stackCallsToSkip+2)

	// Assegna direttamente il puntatore al payload,
	// dal momento che l'item viene generato e immediatamente loggato
	// senza essere ulteriormente manipolato.
	l.muPayload.RLock()
	item.Payload = l.payload
	l.muPayload.RUnlock()

	if l.initItemF != nil {
		l.initItemF(item)
	}

	globalDispatcher.Dispatch(item)
}
