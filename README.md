# vbf (Very Best Framework)
**v**ery **b**est **f**ramework ✨

A set of functions which make it easier to work with the go standard http library. That's it. 💣

## Installation
```bash
go get github.com/Phillip-England/vbf@v0.1.0
```

## Quickstart
the quickest way to get a server up
```go
package main

import (
	"html/template"
	"net/http"

	"github.com/Phillip-England/vbf"
)

const KEY_TEMPLATES = "TEMPLATES"

func main() {

	mux, gCtx := vbf.VeryBestFramework()

	templates, err := vbf.ParseTemplates("./templates", nil)
	if err != nil {
		panic(err)
	}
	vbf.SetGlobalContext(gCtx, KEY_TEMPLATES, templates)

	vbf.HandleStaticFiles(mux)
	vbf.HandleFavicon(mux)

	vbf.AddRoute("GET /", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		templates, _ := vbf.GetContext(KEY_TEMPLATES, r).(*template.Template)
		vbf.WriteHTML(w, "<h1>Hello, World!</h1>")
	}, vbf.MwLogger)

	err = vbf.Serve(mux, "8080")
	if err != nil {
		panic(err)
	}

}
```

## Global Context

share context values from the outside of the application with the inside middleware/handlers
```go
func main() {
	mux, gCtx := vbf.VeryBestFramework()
  vbf.SetGlobalContext(gCtx, "KEY", "<h1>Hello, Context!</h1>") // <--- string
	vbf.AddRoute("GET /", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
    val, _ := vbf.GetContext("KEY", r).(string) // <--- convert back to string
		vbf.WriteHTML(w, val)
	}, vbf.MwLogger)
	err := vbf.Serve(mux, "8080")
	if err != nil {
		panic(err)
	}
}
```

## Middleware

copy and paste the middleware skelenton
```go
func NewMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// logic before request
		next.ServeHTTP(w, r)
        // logic after request
	})
}
```

use the middleware in a handler (the last middleware in the chain will be called first)
```go
func main() {
    // --snip

	vbf.AddRoute("GET /", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
        val, _ := vbf.GetContext("KEY", r).(string) // <--- convert back to string
		vbf.WriteHTML(w, val)
	}, vbf.MwLogger, NewMiddleware) // <--- middlware chained

    // --snip
}
```
