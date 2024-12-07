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

	"github.com/PuerkitoBio/goquery"
	"github.com/a-h/templ"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

//=====================================
// INIT
//=====================================

// gives you a few things you need to get an app up and running
func VeryBestFramework() (mux *http.ServeMux, gCtx map[string]any) {
	mux = http.NewServeMux()
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
	contentTypes := map[string]string{
		".html":  "text/html",
		".css":   "text/css",
		".js":    "application/javascript",
		".png":   "image/png",
		".jpg":   "image/jpeg",
		".jpeg":  "image/jpeg",
		".gif":   "image/gif",
		".svg":   "image/svg+xml",
		".json":  "application/json",
		".xml":   "application/xml",
		".txt":   "text/plain",
		".pdf":   "application/pdf",
		".woff":  "font/woff",
		".woff2": "font/woff2",
		".ttf":   "font/ttf",
		".eot":   "application/vnd.ms-fontobject",
		".ico":   "image/x-icon",
		".zip":   "application/zip",
		".tar":   "application/x-tar",
		".gz":    "application/gzip",
	}

	mux.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		chain(func(w http.ResponseWriter, r *http.Request) {
			filePath := r.URL.Path[len("/static/"):]
			fullPath := filepath.Join(".", "static", filePath)
			file, err := os.Open(fullPath)
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			defer file.Close()
			ext := filepath.Ext(filePath)
			contentType, found := contentTypes[ext]
			if !found {
				contentType = "application/octet-stream" // Default content type
			}
			w.Header().Set("Content-Type", contentType)
			http.ServeContent(w, r, filePath, time.Now(), file)
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

// responds from a handler with a string while setting the appropriate headers
func WriteHTML(w http.ResponseWriter, content string) {
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(content))
}

// responds from a handler as plain text while setting the appropriate headers
func WriteString(w http.ResponseWriter, content string) {
	w.Header().Add("Content-Type", "text/plain")
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

// ensures an html string is safe
func SafeString(component string, args ...any) string {
	return template.HTMLEscapeString(fmt.Sprintf(component, args...))
}

// takes a .md file and converts the file to HTML using goldmark, also handlers coloring code blocks
// WARNING: any HTML entities within your markdown content will be loaded AS IS and will not be escaped
// this means this func needs to handled with caution
func LoadMarkdown(mdPath string, theme string) (string, error) {
	if len(mdPath) == 0 {
		fmt.Println("_md elements require a valid path")
	}
	firstChar := string(mdPath[0])
	if firstChar != "." {
		mdPath = "." + mdPath
	}
	mdFileContent, err := os.ReadFile(mdPath)
	if err != nil {
		return "", err
	}
	md := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle(theme),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithHardWraps(),
			goldmarkhtml.WithXHTML(),
			goldmarkhtml.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	err = md.Convert([]byte(mdFileContent), &buf)
	if err != nil {
		return "", err
	}
	str := buf.String()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))
	if err != nil {
		return "", err
	}
	doc.Find("*").Each(func(i int, inner *goquery.Selection) {
		nodeName := goquery.NodeName(inner)
		currentStyle, _ := inner.Attr("style")
		switch nodeName {
		case "pre":
			inner.SetAttr("class", "custom-scroll")
			inner.SetAttr("style", currentStyle+"padding: 1rem; font-size: 0.875rem; overflow-x: auto; border-radius: 0.25rem; margin-bottom: 1rem;")
		case "h1":
			inner.SetAttr("style", currentStyle+"font-weight: bold; font-size: 1.875rem; padding-bottom: 1rem;")
		case "h2":
			inner.SetAttr("style", currentStyle+"font-size: 1.5rem; font-weight: bold; padding-bottom: 1rem; padding-top: 0.5rem; border-top-width: 1px; border-top-style: solid; border-color: #1f2937; padding-top: 1rem;")
		case "h3":
			inner.SetAttr("style", currentStyle+"font-size: 1.25rem; font-weight: bold; margin-top: 1.5rem; margin-bottom: 1rem;")
		case "p":
			inner.SetAttr("style", currentStyle+"font-size: 0.875rem; line-height: 1.5; margin-bottom: 1rem;")
		case "ul":
			inner.SetAttr("style", currentStyle+"padding-left: 1.5rem; margin-bottom: 1rem; list-style-type: disc;")
		case "ol":
			inner.SetAttr("style", currentStyle+"padding-left: 1.5rem; margin-bottom: 1rem; list-style-type: decimal;")
		case "li":
			inner.SetAttr("style", currentStyle+"margin-bottom: 0.5rem;")
		case "blockquote":
			inner.SetAttr("style", currentStyle+"margin-left: 1rem; padding-left: 1rem; border-left: 4px solid #ccc; font-style: italic; color: #555;")
		case "code":
			parent := inner.Parent()
			if goquery.NodeName(parent) == "pre" {
				return
			}
			inner.SetAttr("style", currentStyle+"font-family: monospace; background-color: #1f2937; padding: 0.25rem 0.5rem; border-radius: 0.25rem;")
		case "hr":
			inner.SetAttr("style", currentStyle+"border: none; border-top: 1px solid #ccc; margin: 2rem 0;")
		case "a":
			inner.SetAttr("style", currentStyle+"color: #007BFF; text-decoration: none;")
		case "img":
			inner.SetAttr("style", currentStyle+"max-width: 100%; height: auto; border-radius: 0.25rem; margin: 1rem 0;")
		}
	})
	modifiedHTML, err := doc.Html()
	if err != nil {
		return "", err
	}
	return modifiedHTML, nil
}

// generate a random string for your random-string purposes
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
func ParseTemplates(path string, funcMap template.FuncMap) (*template.Template, error) {
	templates := template.New("").Funcs(funcMap)
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
