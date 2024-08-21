# vbf
**v**ery **b**est **f**ramework, or vbf for short. ✨

## What is vbf?
vbf aims to make it as easy as possible to write web servers in go, no matter the cost. 💣

## Quickstart

This snippet will:

1. start an http server on `localhost:8080`

2. handle serving the favicon.ico at `./favicon.ico` and static files found at `./static`

3. response with a "hello world" string if you ping `localhost:8080/` using:

```bash
curl localhost:8080/
```

4. logs the request in the console using the `vbf.Logger` middleware
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