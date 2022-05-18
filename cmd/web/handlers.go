package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/thikade/webgo/ui"
)

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

// ********************************
// home
// ********************************
func home(w http.ResponseWriter, r *http.Request) {

	// Initialize a slice containing the paths to the two files. It's important
	// to note that the file containing our base template must be the *first*
	// file in the slice.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",
	}

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of file
	// paths as a variadic parameter?
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Use the ExecuteTemplate() method to write the content of the "base"
	// template as the response body.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

// ********************************
// get a customer id and handle response
// ********************************
func customer(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	id := v["id"]
	if id == "alex" {
		w.Write([]byte(fmt.Sprintf("Found customer: %s\n", id)))
	} else {
		http.Error(w, fmt.Sprintf("Customer %s not found", id), http.StatusTeapot)
	}
}

// ********************************
// render fixed template from embed.FS
// ********************************
func renderEmbeddedFile(w http.ResponseWriter, r *http.Request) {
	templateFile_demo1, err := ui.EfsFiles.ReadFile("html/embedded.tpl")
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("404: template not found: %s", "html/embedded.tpl"), http.StatusNotFound)
		return
	}
	t_demo1 := template.Must(template.New("table").Parse(string(templateFile_demo1)))

	buffer := bytes.Buffer{}
	err = t_demo1.Execute(&buffer, pets)
	if err != nil {
		panic(err)
	}
	w.Write(buffer.Bytes())
}

// ********************************
// dynamically load file: template_x from local filesystem
// ********************************
func renderTemplate(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	id := v["tpl"]

	myT, err := template.ParseGlob(fmt.Sprintf("ui/html/external_%s.tpl", id))
	if err != nil {
		http.Error(w, fmt.Sprintf("404: template not found: external_%s.tpl", id), http.StatusNotFound)
		return
	}
	if err := myT.Execute(w, pets); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong in Template rendering ", http.StatusInternalServerError)
	}
}

// ********************************
// GET searchForm: display form
// ********************************
func search_GET(w http.ResponseWriter, r *http.Request) {
	render(w, "ui/html/search.html", nil)
}

// ********************************
// POST search: execute search
// ********************************
func search_POST(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate form
	search := &SearchObj{
		Days:  r.PostFormValue("days"),
		Token: r.PostFormValue("token"),
	}
	// log.Println("DEBUG: PRE Validation")
	if search.Validate() == false {
		//return
	}

	search.ExecuteSearch()
	// render results template
	render(w, "ui/html/search.html", search)

}
