package main

import(
	"net/http"
	"log"
	"time"
	"io"
	"strings"
)

type Handler interface {
	http.Handler
	Setup(prefix string)
}

func NewHandler(fcgi_config *FcgiConfig) *http.ServeMux {
	var mux *http.ServeMux
	var prefix string

	mux = http.NewServeMux()

	prefix = fcgi_config.Prefix

	Loop:
	for _, page := range fcgi_config.Page {
		var handler Handler
		page.Type = strings.ToUpper(page.Type)

		switch(page.Type) {
			case PageTypeGeneric:
				handler = &PageGeneric {
					NewPage(page)
				}
				break
			case PageTypeGallery:
				handler = &PageGallery {
					NewPage(page)
				}
				break
			default:
				log.Print("Failed to find a match for page type " + page.Type)
				continue Loop
		}

		handler.Setup(prefix)
		mux.Handle(page.Path, handler)
	}

	return mux
}
