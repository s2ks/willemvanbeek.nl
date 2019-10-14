package main

import(
	"net/http"
	"bytes"
	"time"
	"html/template"
	"log"
	"io"
)

type GenericHandler struct {
	Page *GenericPage

	Path string
	Content string

	Display bool

	GT *GenericTemplate
}

type GenericTemplate struct {
	Prefix string
	Template *template.Template

	LastExec time.Time
	ExecInterval *time.Duration
	LastError error
}

type GenericTemplateData struct {
	Path string
	Title string
	Name string
	Content string //TODO add support for Content field
}

func (h *GenericHandler) GenericTemplateExec(filepath string) {
	var err error
	var buf bytes.Buffer

	data := GenericTemplateData {
		h.Page.Path,
		h.Page.Title,
		h.Page.Name,
		"",
	}

	h.GT.Prefix = filepath

	for _, tmpl := range h.Page.Template {
		_, err = h.GT.Template.ParseFiles(filepath + tmpl.File)

		if err != nil {
			h.GT.LastError = err
			log.Print(err)
			return
		}
	}

	for _, tmpl := range h.Page.Template {
		err = h.GT.Template.ExecuteTemplate(&buf, tmpl.Name, data)

		if err != nil {
			h.GT.LastError = err
			log.Print(err)
			return
		}
	}

	h.GT.LastExec = time.Now()
	h.Content = buf.String()
}

func (h *GenericHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if time.Now().Sub(h.GT.LastExec).Seconds() > h.GT.ExecInterval.Seconds() {
		h.GenericTemplateExec(h.GT.Prefix)
	}

	//Path should be an exact match
	if(h.Path != r.URL.Path) {
		log.Print(r.URL.Path + " Not found (404)")
		http.NotFound(w, r)
		return
	}

	//Act as if page does not exist
	if(h.Display == false) {
		http.NotFound(w, r)
		return
	}

	//If content is empty something went wrong
	if(h.Content == "") {
		log.Print("Error displaying " + r.URL.Path)
		log.Print(h.GT.LastError)
		http.Error(w, "Server error", 500)
		return
	}

	io.WriteString(w, h.Content)
}
