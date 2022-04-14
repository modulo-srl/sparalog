package env

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// GoroutineID returns goroutine ID.
// Taken from https://blog.sgmansfield.com/2015/12/goroutine-ids/
func GoroutineID() string {
	b := make([]byte, 64)

	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]

	//n, _ := strconv.ParseUint(string(b), 10, 64)
	return string(b)
}

// StackTrace returns the stack trace, skipping the top most calls.
func StackTrace(skip int) string {
	const maxStackLength = 50
	stackBuf := make([]uintptr, maxStackLength)
	length := runtime.Callers(skip+2, stackBuf[:])
	stack := stackBuf[:length]

	trace := "STACKTRACE: goroutine #" + GoroutineID()

	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			trace = trace + fmt.Sprintf("\n%s\n\t%s:%d", frame.Function, frame.File, frame.Line)
		}
		if !more {
			break
		}
	}
	return trace
}

// Device returns the device parameters.
func Device() (program, host string) {
	if progName != "" {
		program = progName + " "
	}

	host = fmt.Sprintf(
		"%s %s.%s",
		hostname, runtime.GOOS, runtime.GOARCH,
	)

	return
}

// Runtime returns the runtime parameters.
func Runtime() string {
	return fmt.Sprintf(
		"%s, cpus: %v, maxprocs: %v, goroutines: %v, cgocalls: %v",
		runtime.Version(), runtime.NumCPU(), runtime.GOMAXPROCS(0),
		runtime.NumGoroutine(), runtime.NumCgoCall(),
	)
}

var progName string
var hostname string

func init() {
	progName = os.Args[0]
	hostname, _ = os.Hostname()
}
