package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

type PageContent struct {
	Raw   string
	Error error
}

type Page struct {
	Path  string //url to handle
	Title string //page title
	Name  string //page name

	Display bool

	Type string //TODO remove

	Files []FileTemplate

	ContentChannel chan *PageContent

	Content PageContent
}

func NewPage(data *PageJson) (page *Page) {
	page = new(Page)

	page.New(data)

	return
}

func (p *Page) New(data *PageJson) {

	p.Path = data.Path
	p.Title = data.Title
	p.Name = data.Name
	p.Display = data.Display
	p.Type = strings.ToUpper(data.Type)

	p.ContentChannel = make(chan *PageContent)

	p.Files = make([]FileTemplate, len(data.Template))

	for i, tmpl := range data.Template {
		p.Files[i] = *(NewFileTemplate(&tmpl))
	}
}

func (p *Page) SetContent(content *PageContent) bool {

	p.Content.Raw = content.Raw
	p.Content.Error = content.Error

	return content.Error == nil
}

func (p *Page) GetContent() string {
	return p.Content.Raw
}

func (p *Page) HasContent() bool {
	return p.Content.Error == nil
}

func (p *Page) RegisterTemplateForExec(prefix string, data interface{}, templ PageTemplate) {
	templ.RegisterForExec(prefix, data, p.Files, p.ContentChannel)
}

func (p *Page) Serve(w http.ResponseWriter, r *http.Request) {
	if p.Path != r.URL.Path {
		log.Print(r.URL.Path + " Not found 404")
		http.NotFound(w, r)
		return
	}

	select {
	case c := <-p.ContentChannel:
		p.SetContent(c)
	default:
	}

	if p.Display == false {
		log.Print("Access denied to " + r.URL.Path)
		http.NotFound(w, r)
		return
	}

	if p.HasContent() == false {
		log.Print("Error displaying " + r.URL.Path)
		log.Print(p.Content.Error) //FIXME PageT undefined
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	io.WriteString(w, p.GetContent())
}
