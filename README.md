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
4. log the request in the console using the `vbf.Logger` middleware
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
    vbf.Add(mux, "GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}, vbf.Logger)
    err := vbf.Serve(mux, "8080")
	if err != nil {
		panic(err)
	}
}
```

## Skeletons

handler
```go
vbf.Add(mux, "GET /", func(w http.ResponseWriter, r *http.Request) {

}, vbf.Logger)
```

middleware
```go
func _(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// logic before
		next.ServeHTTP(w, r)
		// logic after request
	})
}
```

## Handlers


### Handler Skeleton
the skeleton for a handler
```go
vbf.Add(mux, "GET /", func(w http.ResponseWriter, r *http.Request) {

}, vbf.Logger)
```

take note, this line is where we chain our middleware
```go
}, vbf.Logger)
```

we could do something like this
```go
}, vbf.Logger, AnotherMiddleware, SomeOtherMiddleware) // ect..
```

### 404 Page

on the `GET /` or `POST /` (or other methods), you can address 404s like so
```go
vbf.Add(mux, "GET /", func(w http.ResponseWriter, r *http.Request) {
    if r.Path.URL != "/" {
        w.Write([]byte("<h1>404 Not Found</h1>"))
        return
    }
    w.Write([]byte("<h1>Hello, World!</h1>"))
}, vbf.Logger)
```

on other routes, just return your response
```go
vbf.Add(mux, "GET /about", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("<h1>Hello, About Page!</h1>"))
}, vbf.Logger)
```

