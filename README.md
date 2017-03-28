# midsimple
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