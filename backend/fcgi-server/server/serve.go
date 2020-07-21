package server

import (
	"fmt"
	"log"
	"net/http"
)

func LogRequest(r *http.Request) {
	log.Print(fmt.Sprintf("\tHTTP Method: %s", r.Method))
	log.Print(fmt.Sprintf("\tProtocol version: %s", r.Proto))
	log.Print(fmt.Sprintf("\tHeader: %+v", r.Header))
	log.Print(fmt.Sprintf("\tClient: %s", r.RemoteAddr))
	log.Print(fmt.Sprintf("\tContent length: %d", r.ContentLength))

	body := make([]byte, 512)

	n, err := r.Body.Read(body)

	if err != nil || n == 0 {
		return
	}

	log.Print(fmt.Sprintf("\tBody (first 512 bytes): %s", string(body)))

}

func InternalServerError(w http.ResponseWriter) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func (h *Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Path != r.URL.Path {
		http.NotFound(w, r)
		log.Print(fmt.Sprintf("%s Not found (404)", r.URL.Path))
		LogRequest(r)
		return
	}

	h.handler.ServeHTTP(w, r, h)
}

func ServeHTTP(w http.ResponseWriter, r *http.Request, h *Handle) {
	if err := h.GetErr(); err != nil {
		InternalServerError(w)
		log.Print(fmt.Sprintf("Error while attempting to serve \"%s\" - %s", r.URL.Path, err))
		LogRequest(r)
		return
	}

	c, err := h.Content()

	if err != nil {
		InternalServerError(w)
		log.Print(fmt.Sprintf("Error while fetching content for \"%s\" - %s", r.URL.Path, err))
		LogRequest(r)
		return
	}

	_, err = w.Write(c)

	if err != nil {
		InternalServerError(w)
		log.Print(fmt.Sprintf("Error while writing a response for \"%s\" - %s", r.URL.Path, err))
		LogRequest(r)
		return
	}
}
