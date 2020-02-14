package main

import (
	"log"
	"net/http"
	"strings"

	"willemvanbeek.nl/backend/config"
	"willemvanbeek.nl/backend/template"
)

type Handler interface {
	http.Handler //defines ServeHTTP
	Setup(webroot string) error
	New(page *config.PageXml) Handler
	Execute(file, id string, post []TemplatePost)
}

type Handle struct {
	Path string

	Channel chan []byte
	Content []byte
}


var serveMux *http.ServeMux
func GetServeMux() *http.ServeMux {
	if serveMux == nil {
		serveMux = http.NewServeMux()
	}

	return serveMux
}

var handles []Handle
func Register(path string, h Handler) {
	mux := GetServeMux()
	c := make(chan []byte)

	if handles == nil {
		handles = make([]Handle, 0)
	}

	handles = append(handles, Handle{Path: path, Channel: c})

	h.Setup()

	RegisterForExec(path, h, Settings.ExecInterval, c)

	mux.Handle(path, h)
	log.Print("Registered handler for " + path)
}

func NewHandler(conf *config.XmlConf) *http.ServeMux {
	var mux *http.ServeMux
	var webroot string

	mux = http.NewServeMux()

	/* location of the html template files to execute */
	webroot = conf.System.WebRoot

Loop:
	for _, page := range conf.Page {
		var handler Handler
		page.Type = strings.ToUpper(page.Type)

		switch page.Type {
		case PageTypeGeneric:
			var h *PageGeneric
			handler = h.New(&page)

			break
		case PageTypeGallery:
			var h *PageGallery
			handler = h.New(&page)
			break
		default:
			log.Print("Failed to find a match for page type " + page.Type)
			continue Loop
		}

		err := handler.Setup(webroot)

		if err != nil {
			log.Print("Error registering handler for " + page.Path)
			log.Print(err)
		} else {
			mux.Handle(page.Path, handler)
			log.Print("Registered handler for " + page.Path)
		}
	}

	return mux
}
