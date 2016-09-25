# Buflog

A buffer log library in golang. It's easy to record and check logs in multi goroutine program by assign eash log an ID. Some codes are based on Logrus.

## Features:

* Buffer and group logs in one goroutine, define logid at will.
* Colorful logs in terminal.
* Compatibility to golang original log package.

## Getting started

```golang
package main

import (
    log "github.com/bitxel/buflog"
    "os"
    "time"
)

func main() {
    log.SetOutput(os.Stdout)
    log.SetBufTimeout(time.Second)
    log.SetLevel(log.InfoLevel)
    log.ID("logid1").Info("this is info level log")
    log.ID("logid1").Warnln("this is warning level log")
    log.ID("logid2").Errorf("Errorf func. err: %s", "error log")
    log.ID("logid2").Debug("this is debug level log, won't show.")
    time.Sleep(time.Second * 2)
}
```

## Screenshot

## TODO

* Check whether in terminal.
* Add delimiters to seperate different group, omit delimiter if same group.

