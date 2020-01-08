package main

import (
	"log"
	"net/http"
	"strings"
)

type Handler interface {
	http.Handler //defines ServeHTTP
	Setup(srvPath string) error
	New(page *PageJson) Handler
}

func NewHandler(fcgiConfig *FcgiConfig) *http.ServeMux {
	var mux *http.ServeMux
	var srvPath string

	mux = http.NewServeMux()

	/* location of the html template files to execute */
	srvPath = fcgiConfig.System.SrvPath

Loop:
	for _, page := range fcgiConfig.Page {
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

		err := handler.Setup(srvPath)

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
