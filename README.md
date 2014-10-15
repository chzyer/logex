GoLogex
=======
[![Build Status](https://travis-ci.org/chzyer/gologex.svg?branch=master)](https://travis-ci.org/chzyer/gologex)
[![GoDoc](https://godoc.org/github.com/chzyer/gologex?status.svg)](https://godoc.org/github.com/chzyer/gologex)


An golang log lib, supports tracking and level, wrap by standard log lib

How To Get
=======
shell
```
go get github.com/chzyer/gologex
```

source code
```{go}
import "github.com/chzyer/gologex" // package name is logex

func main() {
  logex.Info("Hello!")
}
```

Level
=======

```{go}
import "github.com/chzyer/gologex"

func main() {
  logex.Println("")
  logex.Debug("debug staff.") // Only show if has an "DEBUG" named env variable(whatever value).
  logex.Info("info")
  logex.Warn("")
  logex.Fatal("") // also trigger exec "os.Exit(1)"
  logex.Error(err) // print error
  logex.Struct(obj) // print objs follow such layout "%T(%+v)"
  logex.Pretty(obj) // print objs as JSON-style, more readable and hide non-publish properties, just JSON
}
```

Extendability
======

source code
```{go}
type MyStruct struct {
  BiteMe bool
}
```

may change to

```{go}
type MyStruct struct {
  BiteMe bool
  logex.Logger // just this
}

func main() {
  ms := new(MyStruct)
  ms.Info("woo!")
}
```

Runtime Tracking
======
All log will attach theirs stack info. Stack Info will shown by an layout, `{packageName}.{FuncName}:{FileName}:{FileLine}`

```{go}
package main

import "github.com/chzyer/gologex"

func test() {
	logex.Pretty("hello")
}

func main() {
	test()
}
```

response
```
2014/10/10 15:17:14 [main.test:testlog.go:6][PRETTY] "hello"
```

Error Tracking
======
You can track an error if you want.

```{go}
package main

import (
	"github.com/chzyer/gologex"
	"os"
)

func openfile() (*os.File, error) {
	f, err := os.Open("xxx")
	if err != nil {
		err = logex.Track(err)
	}
	return f, err
}

func test() error {
	f, err := openfile()
	if err != nil {
		return logex.Track(err)
	}
	f.Close()
	return nil
}

func main() {
	err := test()
	if err != nil {
		logex.Error(err)
		return
	}
	logex.Info("test success")
}
```


response
```
2014/10/10 15:22:29 [main.main:testlog.go:28][ERROR] [main.openfile:11;main.test:19] open xxx: no such file or directory
```
