# vbf

## Name
**v**ery **b**est **f**ramework, or vbf for short. ✨

## Minimal
We do what we do without changing the Go standards. Get the **** out of my way. 🏃‍♂️ 

## Quickstart

get a server going on `localhost:8080`
```go
package main

import (
    "fmt"
    "github.com/Phillip-England/vbf"
)

func main() {
    
    port := "8080"
    mux := http.NewServeMux()
	
    vbf.HandleFavicon(mux, Logger)
	vbf.HandleStaticFiles(mux, Logger)
	
    Route(mux, "GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}, vbf.Logger) // <---- chain middleware here
	
    fmt.Println("starting server on port " + port + " 💎")
    err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}
```