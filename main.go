package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/customer/{id:[-a-zA-Z_0-9.]+}", func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)
		id := v["id"]
		if id == "alex" {
			w.Write([]byte(fmt.Sprintf("Found customer: %s", id)))
		} else {
			http.Error(w, fmt.Sprintf("Customer %s not found", id), http.StatusNotFound)
		}
	})
	port := 8080
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        router,
	}
	log.Printf("Listening on port: %d\n", port)
	log.Fatal(s.ListenAndServe())

}
