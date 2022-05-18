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

// fsContent holds our static
// // go:embed ui/html/*
// // go:embed public/*
var fsContent embed.FS

// define some data to fill in our templates
var pets = Pets{
	{
		Animal: "Cat",
		Age:    3,
	},
	{
		Animal: "Dog",
		Age:    7,
	},
}

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

	publicFS, err := fs.Sub(fsContent, "public")
	if err != nil {
		log.Fatal(err)
	}

	// get a customer id and handle response
	router.HandleFunc("/customer/{id:[-a-zA-Z_0-9.]+}", customer)
	// renders fixed template from embed.FS
	router.HandleFunc("/render", renderEmbeddedFile)
	// dynamically load file: template_x from FS
	router.HandleFunc("/render/{tpl:[0-9]+}", renderTemplate)
	// GET searchForm: display form
	router.HandleFunc("/search", search_GET).Methods("GET")
	// POST search: execute search
	router.HandleFunc("/search", search_POST).Methods("POST")

	// default: serve from FSm "public"
	router.PathPrefix("/").Handler(http.FileServer(http.FS(publicFS)))

	log.Printf("Listening on port: %d\n", port)
	log.Fatal(s.ListenAndServe())
}
