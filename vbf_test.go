package vbf

import (
	"net/http"
	"testing"
)

func Test_VBF(t *testing.T) {

	mux, gCtx := VeryBestFramework()

	component := SafeString("Hello, %s", "John")

	AddRoute("GET /", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		WriteHTML(w, component)
	}, MwLogger)

	err := Serve(mux, "8080")
	if err != nil {
		panic(err)
	}
}
