package server

import (
	"fmt"
	"log"
	"net/http"

	"willemvanbeek.nl/backend/logger"
)

func LogRequest(r *http.Request) {
	log.Print(fmt.Sprintf("\tHTTP Method: %s", r.Method))
	log.Print(fmt.Sprintf("\tProtocol version: %s", r.Proto))
	log.Print(fmt.Sprintf("\tHeader: %+v", r.Header))
	log.Print(fmt.Sprintf("\tClient: %s", r.RemoteAddr))
	log.Print(fmt.Sprintf("\tContent length: %d", r.ContentLength))
}

func LogBody(r *http.Request, l uint64) {
	body := make([]byte, l)

	n, err := r.Body.Read(body)

	if err != nil || n == 0 {
		return
	}

	log.Printf("Body (%v byte(s)): %s", l, string(body))
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
	logger.Error("Error while serving %s - %v", r.URL.Path, r)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
	logger.Error("%s not found (404)", r.URL.Path)
}
