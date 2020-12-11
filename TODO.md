## TODO
- Readme

### Writers
- If a writer fails to logs, the failure should be logged to the default writer of the level. 
Must to be solved for writers with worker too.

### Tests
- TestGoroutinePanic() uses StartPanicWatcher() that is incompatible with multiple tests.  
The test should to be moved and tested apart in some way.
