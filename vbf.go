package vbf

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//=====================================
// INIT
//=====================================

// gives you a few things you need to get an app up and running
func VeryBestFramework() (*http.ServeMux, map[string]any) {
	mux := http.NewServeMux()
	HandleFavicon(mux)
	HandleStaticFiles(mux)
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
			http.ServeFile(w, r, fullPath)
		}, middleware...).ServeHTTP(w, r)
	})
}

//=====================================
// REQUEST / RESPONSE HELPERS
//=====================================

// responses from a handler with a string while setting the appropriate headers
func WriteHTML(w http.ResponseWriter, content string) {
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(content))
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
