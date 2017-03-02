// Note: '*' = passing by reference. When you wanna pass ref parameter use '&' prefix

package main

// import required packages for a web application
import (
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

// Page structure
type Page struct {
	Title string
	Body  []byte
}

// Everythig start from here
// Config routing (url => handler)
func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	http.ListenAndServe(":8080", nil)
}

// Attach hanlder to url
// get input as a function which take (http.ResponseWriter, *http.Request, string) as paramters
// return http.HandlerFunc
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(resWriter http.ResponseWriter, req *http.Request) {
		matchedPaths := validPath.FindStringSubmatch(req.URL.Path)
		if matchedPaths == nil {
			http.NotFound(resWriter, req) //if url not match the validPath, return NotFound
			return
		}
		fn(resWriter, req, matchedPaths[2]) // revoke handler function by supply http.ResponseWriter, *http.Request, string)
	}
}

// Handler for url View
// it will load View page
func viewHandler(resWriter http.ResponseWriter, req *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		http.Redirect(resWriter, req, "/edit/"+title, http.StatusFound) //if error, it's mean the title is not existing so load page edit to create new file
		return
	}
	renderTemplate(resWriter, "view", page) // if title exist, render page on template
}

// Handler for url Edit
// it will load Edit page
func editHandler(resWriter http.ResponseWriter, req *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title} // if page not exist, create new one
	}
	renderTemplate(resWriter, "edit", page)
}

// Handler for url Save
// it will save data to a file
func saveHandler(resWriter http.ResponseWriter, req *http.Request, title string) {
	body := req.FormValue("body") //So easy to get form data
	page := &Page{Title: title, Body: []byte(body)}
	err := page.save() //save data to file
	if err != nil {
		http.Error(resWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(resWriter, req, "/view/"+title, http.StatusFound) // redirect to page View after saved
}

// Load content of the specific page from file
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename) // read content in text file as page content
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// Excute html template
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(resWriter http.ResponseWriter, tmpl string, page *Page) {
	err := templates.ExecuteTemplate(resWriter, tmpl+".html", page) //Attach model(Page) to view
	if err != nil {
		http.Error(resWriter, err.Error(), http.StatusInternalServerError)
	}
}

// This is a extension method for object Page
// This example is just save somthing in page body to a file
// and return error if any
func (page *Page) save() error {
	filename := page.Title + ".txt"
	return ioutil.WriteFile(filename, page.Body, 0600)
}
