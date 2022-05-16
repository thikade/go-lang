package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// docRoot holds our static
//go:embed public/*
// //go:embed templates/*
var docRoot embed.FS

func main() {
	router := mux.NewRouter()
	port := 8080
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        router,
	}

	publicFS, err := fs.Sub(docRoot, "public")
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/customer/{id:[-a-zA-Z_0-9.]+}", func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)
		id := v["id"]
		if id == "alex" {
			w.Write([]byte(fmt.Sprintf("Found customer: %s", id)))
		} else {
			http.Error(w, fmt.Sprintf("Customer %s not found", id), http.StatusNotFound)
		}
	})

	// default server from FS
	router.PathPrefix("/").Handler(http.FileServer(http.FS(publicFS)))

	log.Printf("Listening on port: %d\n", port)
	log.Fatal(s.ListenAndServe())

}
