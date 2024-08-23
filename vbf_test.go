package vbf

import (
	"net/http"
	"testing"
)

func Test_VBF(t *testing.T) {

	mux, gCtx := VeryBestFramework()

	AddRoute("GET /", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		WriteHTML(w, "<h1>Hello, World</h1>")
	}, MwLogger)

	err := Serve(mux, "8080")
	if err != nil {
		panic(err)
	}
}
