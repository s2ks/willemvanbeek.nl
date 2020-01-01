package main

import(
	"net/http"
	"log"
	"time"
	"io"
	"strings"
)

type Handler interface {
	CheckTime() bool

	LastError() error
	GetHandlerData() HandlerData

	TemplateExec(string) error
}


func 

func (h T) CheckPath() (isPath bool) {

	if h.Path == r.URL.Path {
		isPath = true
	} else {
		isPath = false;
	}

	return
}

func (h T) DoDisplay() (doDisplay bool) {
	doDisplay = h.Display;

	return
}

func (h T) HasContent() (hasContent bool) {

	if h.Content == "" {
		hasContent = false
	} else {
		hasContent = true
	}

	return
}

func (h T) DoRefresh() (doRefresh bool) {
	doRefresh = false

	lastExec := h.T.LastExec;
	interval := h.T.ExecInterval;

	if time.Now().Sub(*lastExec).Seconds() > interval.Seconds() {
		doRefresh = true
	}

	return
}

func (h T) LastError() error {
	return h.T.LastError
}

func NewHandler(fcgi_config *FcgiConfig) *http.ServeMux {
	var mux *http.ServeMux

	mux = http.NewServeMux()

	for _, page := range fcgi_config.Page {
		page.Type = strings.ToUpper(page.Type)

		switch(page.Type) {
			case PageTypeGeneric:
				handler := &PageGeneric {
					NewPage(&page)
				}
				break
			case PageTypeGallery:
				handler := &PageGallery {
					NewPage(&page)
				}
				break
		}

		mux.Handle(page.Path, handler)
	}

	return mux
}

func (h T) HandleServeHTTP(w http.ResponseWriter, r *http.Request) {

	//Path should be an exact match
	if h.CheckPath() == false {
		log.Print(r.URL.Path + " Not found (404)")
		http.NotFound(w, r)
		return
	}

	//Execute the template again every so-often
	if h.DoRefresh() == true {
		h.TemplateExec(h.T.Prefix)
	}

	//Act as if page does not exist
	if h.DoDisplay() == false {
		log.Print("Access denied to " + r.URL.Path)
		http.NotFound(w, r)
		return
	}

	//If content is empty something went wrong
	if h.HasContent() == false {
		log.Print("Error displaying " + r.URL.Path)
		log.Print(h.LastError())
		http.Errror(w, "Internal server error", 500);
		return
	}

	io.WriteString(w, h.Content)
}
