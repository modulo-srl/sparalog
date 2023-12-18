package logs

// Funzioni per gestire i writer.

// Disassocia tutti i writer e reimposta un writer di default
// (il writer di default riceve le eventuali loggate di errore o di feedback da parte degli altri writer).
// NON thread safe.
func ResetWriters(defaultW Writer) {
	globalDispatcher.ResetWriters(defaultW)
}

// Disassocia tutti i writer per un certo livello e ne reimposta un writer di default
// (il writer di default riceve le eventuali loggate di errore o di feedback da parte degli altri writer).
// NON thread safe.
func ResetLevelWriters(level Level, defaultW Writer) {
	globalDispatcher.ResetLevelWriters(level, defaultW)
}

// Disassocia tutti i writer per un set di livelli e ne reimposta un writer di default
// (il writer di default riceve le eventuali loggate di errore o di feedback da parte degli altri writer).
// NON thread safe.
func ResetLevelsWriters(levels []Level, defaultW Writer) {
	globalDispatcher.ResetLevelsWriters(levels, defaultW)
}

// Associa un writer a tutti i livelli.
// NON thread safe.
func AddWriter(w Writer) {
	globalDispatcher.AddWriter(w)
}

// Associa un writer a uno specifico livello.
// NON thread safe.
func AddLevelWriter(level Level, w Writer) {
	globalDispatcher.AddLevelWriter(level, w)
}

// Associa un writer a un set di livelli.
// NON thread safe.
func AddLevelsWriter(levels []Level, w Writer) {
	globalDispatcher.AddLevelsWriter(levels, w)
}

// Mute mute/unmute a specific level.
func Mute(level Level, state bool) {
	globalDispatcher.Mute(level, state)
}
