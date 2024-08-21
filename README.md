# vbf
**v**ery **b**est **f**ramework, or vbf for short. ✨

## Philsophy
We do what we do without changing the go standards. All the code **you'll** write when using vbf could *easily* be migrated over to a vanilla go project. ❤️ 

## Quickstart

This snippet will:

1. Get an http server running on `localhost:8080`. 

2. Handle serving the favicon.ico at `./favicon.ico` and static files found at `./static`.

3. Sends a "hello world" string back if you ping `localhost:8080/`

4. Logs the request in the console using the `vbf.Logger` middleware.
```go
package main

import (
    "fmt"
    "github.com/Phillip-England/vbf"
)

func main() {
    mux := http.NewServeMux()
    vbf.HandleFavicon(mux, Logger)
	vbf.HandleStaticFiles(mux, Logger)
    vbf.Route(mux, "GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}, vbf.Logger)
    err := vbf.Serve(mux, "8080")
	if err != nil {
		panic(err)
	}
}
```