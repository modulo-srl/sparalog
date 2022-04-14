package sparalog

// Item is the single item of log.
type Item interface {
	Context

	Level() Level

	Line() string
	SetLine(string)

	StackTrace() string
	GenerateStackTrace(callsToSkip int)
	SetStackTrace(string)

	ToString(timestamp, level bool) string

	Log()

	// Internal domain functions - should not be used out of this library.

	SetLogger(ItemLogger)
}

// ItemLogger used by Item.SetLogger(), Item.Log().
type ItemLogger interface {
	LogItem(Item)
}

// Context defines the interface for contextualized logging data.
type Context interface {
	SetTag(string, string)
	SetData(string, interface{})
	SetPrefix(format string, tags []string)

	// Internal domain functions - should not be used out of this library.

	AssignContext(Context, bool)

	Tags() map[string]string
	Data() map[string]interface{}

	Prefix() (string, []string)

	RenderPrefix() string
}
