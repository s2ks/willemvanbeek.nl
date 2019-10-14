package main

import(
	"net/http"
	"strings"
)

type PostHandler struct {
	Path string
}

func (h *PostHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Print("Post Handler, got method: " + r.Method)
		return
	}

	if strings.ToUpper(h.Path) != strings.ToUpper(r.URL.Path) {
		log.Print("Post Handler: path mismatch: " + r.URL.Path + " != " + h.Path)
		return;
	}
}
