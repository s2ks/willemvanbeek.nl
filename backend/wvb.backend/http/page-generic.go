package main

import(
	"net/http"
	"bytes"
	"time"
	"html/template"
	"log"
)
type PageGeneric struct {
	Page
}

type GenericData struct {
	Path string
	Title string
	Name string
}

//TODO implement
func (p *PageGeneric) ExecTemplate() {
}

//TODO DRY
func (h *GenericHandler) TemplateExec(filepath string) (err error) {
	var buf bytes.Buffer

	data := GenericTemplateData {
		h.Page.Path,
		h.Page.Title,
		h.Page.Name,
		"",
	}

	err = nil

	h.T.Prefix = filepath
	h.T.LastExec = time.Now()

	for _, tmpl := range h.Page.Template {
		_, err = h.T.Template.ParseFiles(filepath + tmpl.File)

		if err != nil {
			h.T.LastError = err
			log.Print(err)
			return
		}
	}

	for _, tmpl := range h.Page.Template {
		err = h.T.Template.ExecuteTemplate(&buf, tmpl.Name, data)

		if err != nil {
			h.T.LastError = err
			log.Print(err)
			return
		}
	}

	h.Content = buf.String()

	return
}


func (p *PageGeneric) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Serve(w, r)
}
