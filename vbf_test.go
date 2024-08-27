package vbf

import (
	"html/template"
	"net/http"
	"testing"
)

const KEY_TEMPLATES = "TEMPLATES"

func Test_VBF(t *testing.T) {

	mux, gCtx := VeryBestFramework()

	templates, err := ParseTemplates("./templates")
	if err != nil {
		panic(err)
	}

	SetGlobalContext(gCtx, KEY_TEMPLATES, templates)

	AddRoute("GET /", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		templates, _ := GetContext(KEY_TEMPLATES, r).(*template.Template)
		mdContent, _ := LoadMarkdown("./docs/index.md")
		ExecuteTemplate(w, templates, "layout.html", map[string]interface{}{
			"Title":      "very best framework",
			"HeaderText": "vbf",
			"SubText":    "very best framework",
			"MdContent":  template.HTML(mdContent),
		})
	}, MwLogger)

	err = Serve(mux, "8080")
	if err != nil {
		panic(err)
	}
}
