package vbf

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Phillip-England/ffh"
	"github.com/a-h/templ"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/renderer/html"
)

//=====================================
// INIT
//=====================================

// gives you a few things you need to get an app up and running
func VeryBestFramework() (*http.ServeMux, map[string]any) {
	mux := http.NewServeMux()
	return mux, make(map[string]any)
}

//=====================================
// MIDDLEWARE
//=====================================

// used to chain middleware and handlers in the proper sequence
func chain(h http.HandlerFunc, middleware ...func(http.Handler) http.Handler) http.Handler {
	finalHandler := http.Handler(h)
	for _, m := range middleware {
		finalHandler = m(finalHandler)
	}
	return finalHandler
}

// a middleware to test setting the request context
func MwContentHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		next.ServeHTTP(w, r)
	})
}

// a logging middleware which logs out details about the request
func MwLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		endTime := time.Since(startTime)
		fmt.Printf("[%s][%s][%s]\n", r.Method, r.URL.Path, endTime)
	})
}

//=====================================
// PROVIDED HANDLERS
//=====================================

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
			ext := strings.ToLower(filepath.Ext(filePath))
			var contentType string
			switch ext {
			case ".js":
				contentType = "application/javascript"
			case ".css":
				contentType = "text/css"
			default:
				contentType = "application/octet-stream"
			}
			w.Header().Set("Content-Type", contentType)
			http.ServeFile(w, r, fullPath)
		}, middleware...).ServeHTTP(w, r)
	})
}

//=====================================
// REQUEST / RESPONSE HELPERS
//=====================================

// gets a url param
func Param(r *http.Request, paramName string) string {
	return r.URL.Query().Get(paramName)
}

// compares a provided value to a url query param
func ParamIs(r *http.Request, paramName string, valueToCheck string) bool {
	return r.URL.Query().Get(paramName) == valueToCheck
}

// responses from a handler with a string while setting the appropriate headers
func WriteHTML(w http.ResponseWriter, content string) {
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(content))
}

// responds from a handler with a templ component with the appropriate headers
func WriteTempl(w http.ResponseWriter, r *http.Request, component templ.Component) error {
	w.Header().Add("Content-Type", "text/html")
	err := component.Render(r.Context(), w)
	if err != nil {
		return err
	}
	return nil
}

//=====================================
// ROUTING
//=====================================

// adds a new route to your server
func AddRoute(path string, mux *http.ServeMux, globalCtx map[string]any, handler http.HandlerFunc, middleware ...func(http.Handler) http.Handler) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		r = SetContext("GLOBAL", globalCtx, r)
		chain(handler, middleware...).ServeHTTP(w, r)
	})
}

//=====================================
// CONTEXT HELPERS
//=====================================

// used to setup key for setting request context data
type ContextKey string

// retrieves a map to be used as a global context for the app
func MakeGlobalContext() map[string]any {
	return make(map[string]any)
}

// sets a value on the global context
func SetGlobalContext(globalCtx map[string]any, key string, value any) {
	globalCtx[key] = value
}

// to be used inside a middleware or handler to share context data with other middleware/handlers
func SetContext(key string, val any, r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), ContextKey(key), val)
	r = r.WithContext(ctx)
	return r
}

// to be used inside a middleware or handler to get context data set in other middleware/handlers
func GetContext(key string, r *http.Request) any {
	ctxMap, ok := r.Context().Value(ContextKey("GLOBAL")).(map[string]any)
	if ok {
		mapVal := ctxMap[key]
		if mapVal != nil {
			return mapVal
		}
	}
	val := r.Context().Value(ContextKey(key))
	return val
}

//=====================================
// GENERALY UTILITY FUNCS
//=====================================

func SafeString(component string, args ...any) string {
	return template.HTMLEscapeString(fmt.Sprintf(component, args...))
}

func LoadMarkdown(filepath string) (string, error) {
	mdContent, err := ffh.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	md := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(mdContent), &buf); err != nil {
		return "", err
	}
	html := buf.String()
	// giving all <pre> a class so I can access them via JS

	preClassName, err := UtilRandStr(32)
	if err != nil {
		return "", err
	}
	html = strings.ReplaceAll(html, "<pre", fmt.Sprintf(`<pre class="%s"`, preClassName))
	html = strings.ReplaceAll(html, fmt.Sprintf(`<pre class="%s" style="`, preClassName), fmt.Sprintf(`<pre class="%s text-sm flex flex-col relative p-2" style="overflow-x:auto;`, preClassName))
	copyButton := `<div id='copy-button' class='flex items-center cursor-pointer absolute right-0'><svg class="w-6 h-6 text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24"><path fill-rule="evenodd" d="M18 3a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2h-1V9a4 4 0 0 0-4-4h-3a1.99 1.99 0 0 0-1 .267V5a2 2 0 0 1 2-2h7Z" clip-rule="evenodd"/><path fill-rule="evenodd" d="M8 7.054V11H4.2a2 2 0 0 1 .281-.432l2.46-2.87A2 2 0 0 1 8 7.054ZM10 7v4a2 2 0 0 1-2 2H4v6a2 2 0 0 0 2 2h7a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3Z" clip-rule="evenodd"/></svg></div>`
	copyIndicator := `<div id='copy-indicator' class='hidden flex items-center absolute right-0'><svg class="w-6 h-6 text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24"><path fill-rule="evenodd" d="M9 2a1 1 0 0 0-1 1H6a2 2 0 0 0-2 2v15a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V5a2 2 0 0 0-2-2h-2a1 1 0 0 0-1-1H9Zm1 2h4v2h1a1 1 0 1 1 0 2H9a1 1 0 0 1 0-2h1V4Zm5.707 8.707a1 1 0 0 0-1.414-1.414L11 14.586l-1.293-1.293a1 1 0 0 0-1.414 1.414l2 2a1 1 0 0 0 1.414 0l4-4Z" clip-rule="evenodd"/></svg></div>`
	html = strings.ReplaceAll(html, "<code>", `<div class='flex flex-row items-center gap-2 p-4 relative'>`+copyIndicator+copyButton+`</div><code class='whitespace-pre-wrap'>`)
	html = html + `<script>
		(() => {
			let blocks = document.querySelectorAll('.` + preClassName + `')
			for (let i = 0; i < blocks.length; i++) {
				let block = blocks[i]
				let copyButton = block.querySelector('#copy-button')
				let copyIndicator = block.querySelector('#copy-indicator')
				let lines = block.textContent.split('\n')
				let newLines = []
				for (let i2 = 0; i2 < lines.length; i2++) {
					let line = lines[i2]
					for (let i3 = 0; i3 < line.length; i3++) {
						let char = line[i3]
						let num = parseInt(char, 10)
						if (Number.isInteger(num)) {
							continue
						} 
						newLines.push(line.slice(i2))
						break
					}
				}
				let htmlData = newLines.join("\n")
				let firstTenChars = htmlData.slice(0, 10)
				let parts = firstTenChars.split(".")
				let firstPart = parts[0]
				let firstPartFirstChar = parts[0][0]
				let isNum = parseInt(firstPartFirstChar, 10)
				if (firstPart && Number.isInteger(isNum)) {
					htmlData = htmlData.slice(htmlData.indexOf(firstPartFirstChar)+1, htmlData.length)
				}
				block.addEventListener('click', (e) => {
					copyButton.classList.add('hidden')
					copyIndicator.classList.remove('hidden')
					navigator.clipboard.writeText(htmlData)
					setTimeout(() => {
						copyButton.classList.remove('hidden')
						copyIndicator.classList.add('hidden')
					}, 1000)
				})
			} 
		})()
	</script>`
	return html, nil
}

func UtilRandStr(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i := 0; i < n; i++ {
		bytes[i] = letters[int(bytes[i])%len(letters)]
	}

	return string(bytes), nil
}

//=====================================
// TEMPLATING HELPERS
//=====================================

// executes an html templates while writing the appropriate headers
func ExecuteTemplate(w http.ResponseWriter, templates *template.Template, filepath string, data any) error {
	w.Header().Add("Content-Type", "text/html")
	err := templates.ExecuteTemplate(w, filepath, data)
	if err != nil {
		return err
	}
	return nil
}

// parse all the templates found at the provided path
func ParseTemplates(path string) (*template.Template, error) {
	templates := template.New("")
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err := templates.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// TemplateToString renders a template with the provided data and returns the result as a string.
func TemplateToString(templates *template.Template, filepath string, data any) (string, error) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, filepath, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

//=====================================
// SERVING
//=====================================

// serves the application at the given port
func Serve(mux *http.ServeMux, port string) error {
	fmt.Println("starting server on port " + port + " 🚀")
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		return err
	}
	return nil
}
