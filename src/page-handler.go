package main

import(
	"net/http"
	"html/template"
	"log"
	"io"
	"bytes"
	"time"
)

type WvbHandler struct {
	Index int
	Prefix string
	Page WvbPage

	Tmpl *template.Template
	Exec bytes.Buffer
	LastExec *time.Time
}

type WvbData struct {
	Path string
	Title string
	Name string

	Content [][]string
}

var handlers []*WvbHandler

func (h *WvbHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Page.Path != r.URL.Path {
		log.Print(r.URL.Path + " not found")
		http.NotFound(w, r)
		return
	}
	if h.Page.Display == false {
		log.Print(r.URL.Path + " Display = false")
		http.NotFound(w, r)
		return
	}

	log.Print(h.Exec.String())
	io.WriteString(w, h.Exec.String())

	//t, err = template.ParseFiles(h.Files)
	/*

	if err {
		log.Print(err)
		return http.NotFound(w, r)
	}

	for i, template := range h.Templates {
		err = t.ExecuteTemplate(w, template, h)
		if err {
			log.Print(err)
			return http.NotFound(w, r)
		}
	}
	*/
}

func (wvb *WvbHandler) wvb_template_exec(prefix string) {
	var err error

	wvb.Exec.Reset()

	data := WvbData {
		wvb.Page.Path,
		wvb.Page.Title,
		wvb.Page.Name,
		nil,
	}

	for _, tmpl := range wvb.Page.Template {
		data.Content = tmpl.Content

		_, err = wvb.T.ParseFiles(prefix + tmpl.File)

		if err != nil {
			log.Print(err)
			continue
		}

		//log.Print(prefix + tmpl.File)
	}

	for _, tmpl := range wvb.Page.Template {
		wvb.Tmpl.ExecuteTemplate(&wvb.Exec, tmpl.Name, data)
	}
}

func wvb_handler_init(wvb_config *WvbConfig) {

	handlers = make([]*WvbHandler, len(wvb_config.Page))

	for i, page := range wvb_config.Page {
		handler := new(WvbHandler)
		handler.Index = i
		handler.Prefix = wvb_config.Prefix
		handler.Page = page
		handler.T = template.New(page.Name)

		handler.wvb_template_exec(wvb_config.Prefix)

		handlers[i] = handler

		http.Handle(page.Path, handler)
		log.Print("Registered handler for " + page.Path)
	}

	/*
	for i, path := range WvbConfig.Paths {
		handler := new(WvbPageHandler)
		handler.Path = path
		handler.Index = i
		handler.Templates = WvbConfig.Templates[i]
		handler.Title = WvbConfig.Title[i]
		handler.Files = WvbConfig.Files[i]
		handler.Content = WvbConfig.Content[i]

		http.Handle(path, handler)
	}
	*/
}
