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

	templateFile_demo1, err := fsContent.ReadFile("templates/demo1.tpl")
	if err != nil {
		log.Fatal(err)
	}
	t_demo1 := template.Must(template.New("table").Parse(string(templateFile_demo1)))

	// ********************************
	// get a customer id and handle response
	// ********************************
	router.HandleFunc("/customer/{id:[-a-zA-Z_0-9.]+}", func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)
		id := v["id"]
		if id == "alex" {
			w.Write([]byte(fmt.Sprintf("Found customer: %s", id)))
		} else {
			http.Error(w, fmt.Sprintf("Customer %s not found", id), http.StatusNotFound)
		}
	})

	// ********************************
	// render fixed template from embed.FS
	// ********************************
	router.HandleFunc("/render", func(w http.ResponseWriter, r *http.Request) {

		buffer := bytes.Buffer{}
		err = t_demo1.Execute(&buffer, pets)
		if err != nil {
			panic(err)
		}
		w.Write(buffer.Bytes())

	})

	// ********************************
	// dynamically load file: template_x from FS
	// ********************************
	router.HandleFunc("/render/{tpl:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)
		id := v["tpl"]

		myT, err := template.ParseGlob(fmt.Sprintf("templates/external_%s.tpl", id))
		if err != nil {
			http.Error(w, fmt.Sprintf("404: template not found: templates/external_%s.tpl", id), http.StatusNotFound)
			return
		}

		buffer := bytes.Buffer{}
		err = myT.Execute(&buffer, pets)
		if err != nil {
			panic(err)
		}
		w.Write(buffer.Bytes())

	})

	// ********************************
	// GET searchForm: display form
	// ********************************
	router.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		render(w, "templates/search.html", nil)
	}).Methods("GET")

	// ********************************
	// POST search: execute search
	// ********************************
	router.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Validate form
		msg := &SearchObj{
			Days:  r.PostFormValue("days"),
			Token: r.PostFormValue("token"),
		}
		// log.Println("DEBUG: PRE Validation")
		if msg.Validate() == false {
			render(w, "templates/search.html", msg)
			return
		}

		// log.Println("DEBUG: Showing results")
		obj := &SearchObj{
			TotalResults: 5,
			Results:      []string{"FOO", "B", "C", "D", "E"},
		}

		// render results template
		render(w, "templates/search.html", obj)

	}).Methods("POST")

	// ********************************
	// default: serve from FS
	// ********************************
	router.PathPrefix("/").Handler(http.FileServer(http.FS(publicFS)))

	log.Printf("Listening on port: %d\n", port)
	log.Fatal(s.ListenAndServe())
}

// renders templates from filesystem (in combination with optional data)
// https://www.alexedwards.net/blog/form-validation-and-processing
func render(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong in Template rendering ", http.StatusInternalServerError)
	}
}
