package server

import (
	"fmt"
	"log"
	"net/http"
)

/* TODO variable body length */
func LogRequest(r *http.Request) {
	log.Print(fmt.Sprintf("\tHTTP Method: %s", r.Method))
	log.Print(fmt.Sprintf("\tProtocol version: %s", r.Proto))
	log.Print(fmt.Sprintf("\tHeader: %+v", r.Header))
	log.Print(fmt.Sprintf("\tClient: %s", r.RemoteAddr))
	log.Print(fmt.Sprintf("\tContent length: %d", r.ContentLength))

	body := make([]byte, 64)

	n, err := r.Body.Read(body)

	if err != nil || n == 0 {
		return
	}

	log.Print(fmt.Sprintf("\tBody (first 64 bytes): %s", string(body)))

}

func InternalServerError(w http.ResponseWriter) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Path != r.URL.Path {
		http.NotFound(w, r)
		log.Print(fmt.Sprintf("%s Not found (404)", r.URL.Path))
		LogRequest(r)
		return
	} else {
		h.ihandler.ServeHTTP(w, r, h)
	}

}
