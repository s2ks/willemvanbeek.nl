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
}

func NewPage(data *PageData) (page *Page) {
	page = new(Page)

	page.New(data)

	return
}

func (p *Page) New(data *PageData) {

	p.Path = data.Path
	p.Title = data.Title
	p.Name = data.Name
	p.Display = data.Display
	p.Type = strings.ToUpper(data.Type)

	p.Files = make([]FileTemplate, len(data.Template))

	for i, tmpl := range data.Template {
		p.FileT[i] = NewFileTemplate(tmpl)
	}
}

func (p *Page) SetContent(content string) bool {

	p.Content = content
	p.HasContent = content != ""

	return p.HasContent
}

func (p *Page) Serve(w http.ResponseWriter, r *http.Request) {
	if p.Path != r.URL.Path {
		log.Print(r.URL.Path + " Not found 404")
		http.NotFound(w, r)
		return
	}

	if p.Display == false {
		log.Print("Access denied to " + r.URL.Path)
		http.NotFound(w, r)
		return
	}

	if p.HasContent == false {
		log.Print("Error displaying " + r.URL.Path)
		log.Print(p.PageT.LastError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}


	io.WriteString(w, p.Content)
}
