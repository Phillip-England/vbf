package vbf

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_Xerus(t *testing.T) {
	port := "8080"
	mux := http.NewServeMux()
	HandleFavicon(mux, Logger)
	HandleStaticFiles(mux, Logger)
	Route(mux, "GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}, Logger)

	fmt.Println("starting server on port " + port + " 💎")
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		panic(err)
	}
}
