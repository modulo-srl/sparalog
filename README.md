# Sparalog

Logging with independent streaming levels.

![logger diagram](/doc/img/logger.svg)

![dispatcher diagram](/doc/img/dispatcher.svg)

## Features

* One logger, multiple writers for every logging level.
* Thread safe.
* Light and tested.
* Logs panics from all goroutines without defer.

## Notes

* Writers internal errors are redirected to the default writer.

---
*Copyright 2020,2023 [Modulo srl](http://www.modulo.srl) - Licensed under the MIT license*
