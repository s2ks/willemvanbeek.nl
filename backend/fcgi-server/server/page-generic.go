package main

import (
	"net/http"
)

type PageGeneric struct {
	Page
	Template PageTemplate
}

type GenericData struct {
	Path  string
	Title string
	Name  string
}

func (p *PageGeneric) New(page *PageJson) Handler {
	if p != nil {
		return p
	}

	return &PageGeneric{
		*(NewPage(page)),
		PageTemplate{},
	}
}

func (p *PageGeneric) Setup(prefix string) error {
	var data GenericData

	data.Path = p.Path
	data.Title = p.Title
	data.Name = p.Name

	p.Template = *(NewPageTemplate())

	p.RegisterTemplateForExec(prefix, data, p.Template)

	return nil
}

func (p *PageGeneric) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	select {
	case c := <-p.ContentChannel:
		p.SetContent(c)
	default:
	}

	p.Serve(w, r)
}
