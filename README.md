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
		// logic before request
		next.ServeHTTP(w, r)
		// logic after request
	})
}
```

## Handlers

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

404 pages can be addressed at `GET /`, `POST /`, `PUT /`, ect
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

## Middleware

the skeleton for middleware
```go
func _(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// logic before request
		next.ServeHTTP(w, r)
		// logic after request
	})
}
```

create a new middleware which sets some data in the request context
```go
func NewMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = vbf.SetCtx("someData", "Hello, Context!", r)
		next.ServeHTTP(w, r)
	})
}
```

create another middleware which gets the context data and logs it to the console 
```go
func AnotherMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := vbf.GetCtx("someData", r).(string)
		next.ServeHTTP(w, r)
	})
}
```

be sure to convert the data back to its original type
```go
r = vbf.SetCtx("someData", "Hello, Context!", r) // context data is a string
val := vbf.GetCtx("someData", r).(string) // annotate type as a string to match
```