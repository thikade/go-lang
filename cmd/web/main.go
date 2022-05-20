package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// define some data to fill in our templates
// Pets tyoe is defined in types.go
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
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	port := flag.String("port", "8080", "HTTP network address")
	flag.Parse()

	outLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// Initialize a new instance of our application struct, containing the
	// dependencies.
	app := &application{
		errorLog: errorLog,
		outLog:   outLog,
	}

	router := mux.NewRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", *port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       errorLog,
		Handler:        router,
	}

	// get a customer id and handle response
	router.HandleFunc("/customer/{id:[-a-zA-Z_0-9.]+}", app.customer)
	// renders fixed template from embed.FS
	router.HandleFunc("/embed", app.renderEmbeddedFile)
	// dynamically load file: template_x from FS
	router.HandleFunc("/render/{tpl:[0-9]+}", app.renderTemplate)
	// GET searchForm: display form
	router.HandleFunc("/search", app.search_GET).Methods("GET")
	// POST search: execute search
	router.HandleFunc("/search", app.search_POST).Methods("POST")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/"))))

	// default: serve from  "./public" folder
	// Create a file server which serves files out of the "./public" directory.
	// router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	router.HandleFunc("/", app.home)

	outLog.Printf("Listening on port: %s\n", *port)
	errorLog.Fatal(s.ListenAndServe())
}
