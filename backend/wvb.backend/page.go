package main

import (
	"net/http"
	"log"
	"strings"
)

type Page struct {
	Path string	//url to handle //TODO rename to Url
	Title string	//page title
	Name string	//page name

	Display bool

	Type string

	FileT []FileTemplate

	Content string
	HasContent bool

	PageT *PageTemplate
}

func NewPage(data *PageData) (page *Page) {
	page = new(Page)
	page.Path = data.Path
	page.Title = data.Title
	page.Name = data.Name
	page.Display = data.Display
	page.Type = strings.ToUpper(data.Type)

	page.Files = make([]FileTemplate, len(data.Template))

	for i, tmpl := range data.Template {
		page.FileT[i] = NewFileTemplate(tmpl)
	}

	page.PageT = NewPageTemplate()

	return
}

func (p *Page) Serve(w http.ResponseWriter, r *http.Request) {
	if p.Path != r.URL.Path {
		log.Print(r.URL.Path + " Not found 404")
		http.NotFound(w, r)
	}

	if p.Display == false {
		log.Print("Access denied to " + r.URL.Path)
		http.NotFound(w, r)
	}

	


	io.WriteString(w, p.Content)
}
