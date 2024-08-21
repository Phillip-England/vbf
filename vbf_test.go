package vbf

import (
	"net/http"
	"testing"
)

func Test_Xerus(t *testing.T) {
	mux := http.NewServeMux()
	HandleFavicon(mux, Logger)
	HandleStaticFiles(mux, Logger)
	Add(mux, "GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}, Logger)
	err := Serve(mux, "8080")
	if err != nil {
		panic(err)
	}
}
