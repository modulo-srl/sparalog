package logs

// Funzioni per la generazione item generici (generati dal logger di default).

// Genera un nuovo item dal logger di default, di livello specifico.
// Eredita una copia del payload dal logger che può essere ulteriormente customizzata.
func NewItem(level Level, args ...any) *Item {
	return defaultLogger.NewItem(level, args...)
}

// Genera un nuovo item dal logger di default, di livello specifico.
// Eredita una copia del payload dal logger che può essere ulteriormente customizzata.
func NewItemf(level Level, format string, args ...any) *Item {
	return defaultLogger.NewItemf(level, format, args...)
}

// Genera un nuovo item dal logger di default, di livello errore.
// Eredita una copia del payload dal logger che può essere ulteriormente customizzata.
func NewErrorItem(err error) *Item {
	return defaultLogger.NewErrorItem(err)
}

// Genera un nuovo item dal logger di default, di livello errore.
// Eredita una copia del payload dal logger che può essere ulteriormente customizzata.
func NewErrorItemf(format string, a ...any) *Item {
	return defaultLogger.NewErrorItemf(format, a...)
}
