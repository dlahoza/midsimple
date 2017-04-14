# midsimple
[![Build Status](https://travis-ci.org/DLag/midsimple.svg?branch=master)](https://travis-ci.org/DLag/midsimple)
[![Go Report Card](https://goreportcard.com/badge/github.com/DLag/midsimple)](https://goreportcard.com/report/github.com/DLag/midsimple)
[![codecov](https://codecov.io/gh/DLag/midsimple/branch/master/graph/badge.svg)](https://codecov.io/gh/DLag/midsimple)

Simple Go HTTP middleware manager

MidSimple delivers idiomatic and convenient way to chain classic `net/http` middlewares and wrap `http.Handler` and `http.HandlerFunc` types with it.

Classic middleware means:
```go
func (http.Handler) http.Handler
```

### Less words more examples:
```go
package main

import (
    "net/http"
    "time"
    
    "github.com/throttled/throttled"
    "github.com/justinas/nosurf"
    "github.com/DLag/midsimple"
)

func timeoutHandler(h http.Handler) http.Handler {
    return http.TimeoutHandler(h, 1*time.Second, "timed out")
}

func myApp(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello world!"))
}

func main() {
    th := throttled.Interval(throttled.PerSec(10), 1, &throttled.VaryBy{Path: true}, 50)
    chain := midsimple.New()
    chain.Use(th.Throttle)
    chain.Use(timeoutHandler)
    chain.Use(nosurf.NewPure)
    http.ListenAndServe(":8000", chain.WrapFunc(myApp))
}
```
### Adding middlewares
```go
// Start from creating new chain
c := midsimple.New()
// Or adding some middleware from start
c = midsimple.New(Middleware0)
// Simple way
c.Use(Middleware1)
// Stacking
c.Use(Middleware2).Use(Middleware3).Use(Middleware3)
// Variadic
c.Use(Middleware4, Middleware5, Middleware6)
// You can clean your chain with Reset and reuse the same memory block
c.Reset()
// Madness :)
c = New(Middleware0, Middleware1).
        Use(Middleware2, Middleware3).
        Reset().
        Use(Middleware4).Use(Middleware5, Middleware6, Middleware7)
```
### Wrapping your handlers
```go
// You can wrap http.Handler or http.HandlerFunc
c.Wrap(h http.Handler)
c.WrapFunc(hf http.HandlerFunc)
// You can revert order of middlewares
c.WrapRevert(h http.Handler)
c.WrapRevertFunc(hf http.HandlerFunc)
// And use wrapped chain for your webserver
http.ListenAndServe(":8080", c.Wrap(yourHandler))
```