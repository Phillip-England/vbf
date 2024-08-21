package vbf

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"
)

type ContextKey string

func chain(h http.HandlerFunc, middleware ...func(http.Handler) http.Handler) http.Handler {
	finalHandler := http.Handler(h)
	for _, m := range middleware {
		finalHandler = m(finalHandler)
	}
	return finalHandler
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		requestID := fmt.Sprintf("%d", startTime.UnixNano())
		ctx := context.WithValue(r.Context(), ContextKey("requestID"), requestID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
		endTime := time.Since(startTime)
		fmt.Printf("[%s][%s][%s]\n", r.Method, r.URL.Path, endTime)
	})
}

func HandleFavicon(mux *http.ServeMux, middleware ...func(http.Handler) http.Handler) {
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		chain(func(w http.ResponseWriter, r *http.Request) {
			filePath := "favicon.ico"
			fullPath := filepath.Join(".", ".", filePath)
			http.ServeFile(w, r, fullPath)
		}, middleware...).ServeHTTP(w, r)
	})
}

func HandleStaticFiles(mux *http.ServeMux, middleware ...func(http.Handler) http.Handler) {
	mux.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		chain(func(w http.ResponseWriter, r *http.Request) {
			filePath := r.URL.Path[len("/static/"):]
			fullPath := filepath.Join(".", "static", filePath)
			http.ServeFile(w, r, fullPath)
		}, middleware...).ServeHTTP(w, r)
	})
}

func Route(mux *http.ServeMux, path string, handler http.HandlerFunc, middleware ...func(http.Handler) http.Handler) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		chain(handler, middleware...).ServeHTTP(w, r)
	})
}

func Serve(mux *http.ServeMux, port string) error {
	fmt.Println("starting server on port " + port + " 💎")
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		return err
	}
	return nil
}
