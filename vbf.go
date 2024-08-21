package vbf

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"
)

// used to setup key for setting request context data
type CtxKey string

// used to chain middleware and handlers in the proper sequence
func chain(h http.HandlerFunc, middleware ...func(http.Handler) http.Handler) http.Handler {
	finalHandler := http.Handler(h)
	for _, m := range middleware {
		finalHandler = m(finalHandler)
	}
	return finalHandler
}

// a middleware to test setting the request context
func mwSetCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = SetCtx("someData", "Hello, World!", r)
		next.ServeHTTP(w, r)
	})
}

// a middleware to test getting the request context
func mwGetCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := GetCtx("someData", r).(string)
		fmt.Println(val)
		next.ServeHTTP(w, r)
	})
}

// a logging middleware which logs out details about the request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		endTime := time.Since(startTime)
		fmt.Printf("[%s][%s][%s]\n", r.Method, r.URL.Path, endTime)
	})
}

// when called, your server will serve a favicon if it is located at `./favicon.ico`
func HandleFavicon(mux *http.ServeMux, middleware ...func(http.Handler) http.Handler) {
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		chain(func(w http.ResponseWriter, r *http.Request) {
			filePath := "favicon.ico"
			fullPath := filepath.Join(".", ".", filePath)
			http.ServeFile(w, r, fullPath)
		}, middleware...).ServeHTTP(w, r)
	})
}

// when called, your server will serve static files if located at `./static`
func HandleStaticFiles(mux *http.ServeMux, middleware ...func(http.Handler) http.Handler) {
	mux.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		chain(func(w http.ResponseWriter, r *http.Request) {
			filePath := r.URL.Path[len("/static/"):]
			fullPath := filepath.Join(".", "static", filePath)
			http.ServeFile(w, r, fullPath)
		}, middleware...).ServeHTTP(w, r)
	})
}

// adds a new route to your server
func Add(mux *http.ServeMux, path string, handler http.HandlerFunc, middleware ...func(http.Handler) http.Handler) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		chain(handler, middleware...).ServeHTTP(w, r)
	})
}

// to be used inside a middleware or handler to share context data with other middleware/handlers
func SetCtx(key string, val any, r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), CtxKey(key), val)
	r = r.WithContext(ctx)
	return r
}

// to be used inside a middleware or handler to get context data set in other middleware/handlers
func GetCtx(key string, r *http.Request) any {
	val := r.Context().Value(CtxKey(key))
	return val
}

// serves the application at the given port
func Serve(mux *http.ServeMux, port string) error {
	fmt.Println("starting server on port " + port + " 💎")
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		return err
	}
	return nil
}
