package main

import (
	"net/http"
)

type PageGallery struct {
	Page
	Template PageTemplate
}

type GalleryData struct {
	Path  string
	Title string
	Name  string

	ImgPath []string
}

func (p *PageGallery) Setup(prefix string) {
	var data GalleryData

	data.Path = p.Path
	data.Title = p.Title
	data.Name = p.Name

	//TODO get images

	p.Template = *(NewPageTemplate())

	p.RegisterTemplateForExec(prefix, data, p.Template)
}

func (p *PageGallery) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Serve(w, r)
}
