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

func (p *PageGeneric) Setup(prefix string) {
	var data GenericData

	data.Path = p.Path
	data.Title = p.Title
	data.Name = p.Name

	p.Template = *(NewPageTemplate())

	p.RegisterTemplateForExec(prefix, data, p.Template)
}

func (p *PageGeneric) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Serve(w, r)
}
