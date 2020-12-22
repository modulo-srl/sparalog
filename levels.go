package sparalog

// Level type.
type Level int

// Logging levels.
const (
	// FatalLevel - Shutdown of the service or application to prevent data loss (or firther data loss). Wake up the SysAdmin!
	// Stack trace enabled by default.
	FatalLevel Level = iota
	// ErrorLevel - Any error which is fatal to the operation, but not the service or application (can't open a required file, missing data, incorrect connection strings, missing services, etc.).
	// SysAdmin should be notified automatically, but doesn't need to be dragged out of bed.
	// Stack trace enabled by default.
	ErrorLevel
	// WarnLevel - Anything that can potentially cause application oddities, but automatically recovered.
	WarnLevel
	// InfoLevel - General operational entries about what's going on inside the service or application.
	// Should be the out-of-the-box level.
	InfoLevel
	// DebugLevel - Usually enabled only when debugging. Very verbose logging.
	// Muted by default.
	DebugLevel
	// TraceLevel - For tracing the code and trying to find one part of a function specifically.
	// Muted by default.
	TraceLevel

	LevelsCount
)

// Levels is a constant of all logging levels.
var Levels = [LevelsCount]Level{
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}

// CriticalLevels lists critical levels.
var CriticalLevels = []Level{FatalLevel, ErrorLevel, WarnLevel}

// DebugLevels lists debugging purpose levels.
var DebugLevels = []Level{DebugLevel, TraceLevel}

// LevelsString is a constant of all logging levels names.
var LevelsString = [LevelsCount]string{
	"fatal", "error", "warning", "info", "debug", "trace",
}

// LevelsIcons is a constant of all logging levels UTF8 icons.
var LevelsIcons = [LevelsCount]string{
	"\xE2\x9D\x8C", "\xE2\x9D\x97", "\xE2\x9A\xA0", "\xE2\x84\xB9", "\xF0\x9F\x90\x9B", "\xF0\x9F\x94\x8E",
}
