package main

import(
	"net/http"
	"bytes"
	"time"
	"html/template"
	"log"
)
//TODO inherit
type GenericHandler struct {
	Page *PageData

	Path string
	Content string

	Display bool

	T *GenericTemplate
}

//TODO inherit
type GenericTemplate struct {
	Prefix string
	Template *template.Template

	LastExec time.Time
	ExecInterval time.Duration
	LastError error
}

//TODO inherit
type GenericTemplateData struct {
	Path string
	Title string
	Name string
	Content string //TODO add support for Content field
}

//TODO DRY
func NewGenericHandler(Page *PageData, Path string, Prefix string, Display bool, ExecInterval *time.Duration) (h *GenericHandler) {
	h = new(GenericHandler)
	h.T = new(GenericTemplate)

	h.Page = Page
	h.Path = Path
	h.T.Prefix = Prefix
	h.Display = Display
	h.T.ExecInterval = *ExecInterval

	h.TemplateExec(Prefix)

	return
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


//TODO DRY
func (h *GenericHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.HandleServeHTTP(w, r)
}
