package logs

// Funzioni per loggare usando il logger di default.

// Logga a livello fatale effetuando anche os.Exit(FatalExitCode)
func Fatal(args ...any) {
	defaultLogger.Fatal(args...)
}

// Logga a livello fatale effetuando anche os.Exit(FatalExitCode)
func Fatalf(format string, args ...any) {
	defaultLogger.Fatalf(format, args...)
}

// Logga a livello errore.
func Error(args ...any) {
	defaultLogger.Error(args...)
}

// Logga a livello errore.
func Errorf(format string, args ...any) {
	defaultLogger.Errorf(format, args...)
}

// Logga a livello warning.
func Warning(args ...any) {
	defaultLogger.Warning(args...)
}

// Logga a livello warning.
func Warningf(format string, args ...any) {
	defaultLogger.Warningf(format, args...)
}

// Logga a livello info.
func Info(args ...any) {
	defaultLogger.Info(args...)
}

// Logga a livello info.
func Infof(format string, args ...any) {
	defaultLogger.Infof(format, args...)
}

// Logga a livello debug.
func Debug(args ...any) {
	defaultLogger.Debug(args...)
}

// Logga a livello debug.
func Debugf(format string, args ...any) {
	defaultLogger.Debugf(format, args...)
}

// Logga un item precedentemente generato.
func LogItem(item *Item) {
	defaultLogger.LogItem(item)
}
