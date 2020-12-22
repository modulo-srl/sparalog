## TODO
- Item.RenderPrefix() with standardized tag for goroutine number.
- Comments for interface methods.
- Update read.me
- Item with prefix and tostring cached (optimizing writers calls to Item's methods)

### TCP writer
- example for auto enable level

### Tests
- TestGoroutinePanic() uses StartPanicWatcher() that is incompatible with multiple tests.  
The test should to be moved and tested apart in some way.
- context_test.go
- fatal_test.go
- file_test.go
