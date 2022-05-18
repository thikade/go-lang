package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

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

	// get a customer id and handle response
	router.HandleFunc("/customer/{id:[-a-zA-Z_0-9.]+}", customer)
	// renders fixed template from embed.FS
	router.HandleFunc("/embed", renderEmbeddedFile)
	// dynamically load file: template_x from FS
	router.HandleFunc("/render/{tpl:[0-9]+}", renderTemplate)
	// GET searchForm: display form
	router.HandleFunc("/search", search_GET).Methods("GET")
	// POST search: execute search
	router.HandleFunc("/search", search_POST).Methods("POST")

	// Create a file server which serves files out of the "./public" directory.
	fileServer := http.FileServer(http.Dir("./public/"))
	// default: serve from  "./public" folder
	router.PathPrefix("/").Handler(fileServer)

	log.Printf("Listening on port: %d\n", port)
	log.Fatal(s.ListenAndServe())
}
