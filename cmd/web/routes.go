package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() *mux.Router {
	// mux := http.NewServeMux()

	// fileServer := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/snippet/view", app.snippetView)
	// mux.HandleFunc("/snippet/create", app.snippetCreate)

	router := mux.NewRouter()
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

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// default: serve from  "./public" folder
	// Create a file server which serves files out of the "./public" directory.
	// router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	router.HandleFunc("/", app.home)

	return router
}
