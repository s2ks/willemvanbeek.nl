package main

import (
	"log"
	"net/http"
	"strings"
)

type Handler interface {
	http.Handler
	Setup(srvPath string)
}

func NewHandler(fcgiConfig *FcgiConfig) *http.ServeMux {
	var mux *http.ServeMux
	var srvPath string

	mux = http.NewServeMux()

	srvPath = fcgiConfig.System.SrvPath

Loop:
	for _, page := range fcgiConfig.Page {
		var handler Handler
		page.Type = strings.ToUpper(page.Type)

		switch page.Type {
		case PageTypeGeneric:
			handler = &PageGeneric{
				*(NewPage(&page)),
				PageTemplate{},
			}
			break
		case PageTypeGallery:
			handler = &PageGallery{
				*(NewPage(&page)),
				PageTemplate{},
			}
			break
		default:
			log.Print("Failed to find a match for page type " + page.Type)
			continue Loop
		}

		handler.Setup(srvPath)
		mux.Handle(page.Path, handler)

		log.Print("Registered handler for " + page.Path)
	}

	return mux
}
