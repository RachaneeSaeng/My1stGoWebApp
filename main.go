package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename) // read content in text file as page content
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):] // slice url path from end off "/view/"
	p, _ := loadPage(title)
	//fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body) // write data to respose writer
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	fmt.Println(t)
	t.Execute(w, p) //render object page on tempplate. {{.Title}} = p.Title
}

func main() {
	http.HandleFunc("/view/", viewHandler) //call function viewHandler on url .../view/
	http.HandleFunc("/edit/", editHandler)
	http.ListenAndServe(":8080", nil) //specifying that it should listen on port 8080 on any interface
}
