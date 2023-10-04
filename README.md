# Sparalog

Logging with independent streaming levels.

![logger diagram](/doc/img/logger.svg)

![dispatcher diagram](/doc/img/dispatcher.svg)

## Features

* One logger, multiple writers for every logging level.
* Thread safe.
* Light and tested.
* Logs panics from all goroutines without defer.

## Usage

```go
import "github.com/modulo-srl/sparalog/logs"

func main() {
    ...
    // Default to stdout/stderror
    logs.Open()
    defer logs.Done()
    ...
    logs.Error("error!")
    logs.Infof("%s", "info!")
    ...
}
```

### Multiple writers

```go
    // New writer to file.
    wf := logs.NewFileWriter("errors.log")

    // Reset the default writer to file for all levels.
    logs.ResetWriters(wf)

    logs.AddLevelsWriter(
        []sparalog.Level{
            sparalog.FatalLevel, sparalog.ErrorLevel, sparalog.WarnLevel,
        },
        ws,
    )
    
    // New Telegram writer.
    wt := logs.NewTelegramWriter("api key", channel id)
    
    // Logs fatals to Telegram too.
    logs.AddLevelWriter(sparalog.FatalLevel, wt)

    logs.Open()
    defer logs.Done()

    ...
```

### Panics watcher

```go
    // Logs fatals to Telegram too.
    wt := logs.NewTelegramWriter("api key", channel id)
    logs.AddLevelWriter(sparalog.FatalLevel, wt, "")

    // Start the watcher.
    // Please note:
    // - Writers, or at least the fatal level writers, must to be set before this calls.
    // - Avoid this call in debugging session!
    logs.StartPanicWatcher()

    logs.Open()

    // Test
    go func() {
        i := 0
        i := 1 / i  // the panic here will be logged
    }()
```

### Misc

```go
    // Enable stack tracke for warning level.
    EnableStacktrace(sparalog.WarnLevel, true)

    ...

    // Mute info, debug and trace levels.
    logs.Mute(sparalog.InfoLevel, true)
    logs.Mute(sparalog.DebugLevel, true)
    logs.Mute(sparalog.TraceLevel, true)

```

### Child loggers

```go
    type module struct {
        log sparalog.Logger
        ...
    }

    func (m *module) init() {
        // This logger will forward to the Default logger writers.
        m.log = logs.NewChildLogger("my module")

        ...

        // Will be logged by Default logger, using "my module" prefix.
        m.log.Error("error!")
    }    
```

## Notes

* Writers internal errors are redirected to the default writer.

---
*Copyright 2020,2023 [Modulo srl](http://www.modulo.srl) - Licensed under the MIT license*
