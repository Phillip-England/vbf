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
	HandleStaticFiles(mux)
	HandleFavicon(mux)

	AddRoute("GET /", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			templates, _ := GetContext(KEY_TEMPLATES, r).(*template.Template)
			mdContent, err := LoadMarkdown("./static/docs/content.md")
			if err != nil {
				WriteString(w, "internal server error")
				return
			}
			ExecuteTemplate(w, templates, "layout.html", map[string]interface{}{
				"Title":      "very best framework",
				"HeaderText": "vbf",
				"SubText":    "very best framework",
				"Content":    template.HTML(mdContent),
			})
		} else {
			WriteString(w, "404 not found")
		}
	}, MwLogger)

	err = Serve(mux, "8080")
	if err != nil {
		panic(err)
	}
}
