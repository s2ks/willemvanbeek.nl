package main

import "net/http"

type WvbPageHandler struct {
	Path string
	Index int

	Templates []string

	Title string
	Files []string
	Content [][]string
}

func (h *WvbPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var t template.Template
	var err error

	if h.Path != r.URL.Path {
		log.Print(r.URL.Path + " not found")
		return http.NotFound(w, r)
	}

	t, err = template.ParseFiles(h.Files)

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
}

func wvb_handler_init() {
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
}
