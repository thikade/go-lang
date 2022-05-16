package main

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

// fsContent holds our static
//go:embed public/*
//go:embed templates/*
var fsContent embed.FS

var pets = []struct {
	Animal string
	Age    int
}{
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

	// https://gowebexamples.com/routes-using-gorilla-mux/
	router.HandleFunc("/customer/{id:[-a-zA-Z_0-9.]+}", func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)
		id := v["id"]
		if id == "alex" {
			w.Write([]byte(fmt.Sprintf("Found customer: %s", id)))
		} else {
			http.Error(w, fmt.Sprintf("Customer %s not found", id), http.StatusNotFound)
		}
	})

	router.HandleFunc("/render", func(w http.ResponseWriter, r *http.Request) {
		templateFile, err := fsContent.ReadFile("templates/demo1.tpl")
		if err != nil {
			log.Fatal(err)
		}
		t := template.Must(template.New("table").Parse(string(templateFile)))

		buffer := bytes.Buffer{}
		err = t.Execute(&buffer, pets)
		if err != nil {
			panic(err)
		}
		w.Write(buffer.Bytes())

	})

	// default server from FS
	router.PathPrefix("/").Handler(http.FileServer(http.FS(publicFS)))

	log.Printf("Listening on port: %d\n", port)
	log.Fatal(s.ListenAndServe())
}
