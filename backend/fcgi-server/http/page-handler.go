package main

//TODO depricate

import(
	"net/http"
	"html/template"
	"log"
	"io"
	//"io/ioutil"
	"bytes"
	"time"
	"strings"
)

type WvbHandler struct {
	Index int
	Prefix string
	Page WvbPage


	Tmpl *template.Template
	Exec bytes.Buffer
	LastExec *time.Time

	LastError error
}

type WvbData struct {
	Path string
	Title string
	Name string

	Action []string
	Method string

	Content [][]string
}

var handlers []*WvbHandler

/*
	Serve h.Exec
*/
func (h *WvbHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if strings.ToUpper(h.Page.Path) != strings.ToUpper(r.URL.Path) {
		log.Print(r.URL.Path + " not found")
		http.NotFound(w, r)
		return
	}
	if h.Page.Display == false {
		log.Print(r.URL.Path + " Display = false")
		http.NotFound(w, r)
		return
	}

	if h.Exec.Len() == 0 {
		log.Print(r.URL.Path + " Could not be displayed")
		log.Print(h.LastError)
		http.Error(w, "Server error", 500)
		return
	}

	io.WriteString(w, h.Exec.String())
}

/*
	Execute all templates in wvb.Tmpl and write result to wvb.Exec
*/
func (wvb *WvbHandler) wvb_template_exec(prefix string) {
	var err error

	data := WvbData {
		wvb.Page.Path,
		wvb.Page.Title,
		wvb.Page.Name,
		wvb.Page.Action,
		wvb.Page.Method,
		nil,
	}

	wvb.Exec.Reset()

	for _, tmpl := range wvb.Page.Template {
		_, err = wvb.Tmpl.ParseFiles(prefix + tmpl.File)

		if err != nil {
			goto exec_err
		}

	}

	for _, tmpl := range wvb.Page.Template {
		data.Content = tmpl.Content
		err = wvb.Tmpl.ExecuteTemplate(&wvb.Exec, tmpl.Name, data)
		if err != nil {
			goto exec_err
		}
	}

	return

exec_err:
	wvb.Exec.Reset()
	log.Print(err)
	wvb.LastError = err
}

/*
	Register handlers for path
*/
func wvb_handler_init(wvb_config *WvbConfig) *http.ServeMux {
	var mux *http.ServeMux

	handlers = make([]*WvbHandler, len(wvb_config.Page))

	mux = http.NewServeMux()


	for i, page := range wvb_config.Page {
		handler := new(WvbHandler)
		handler.Index = i
		handler.Prefix = wvb_config.Prefix
		handler.Page = page
		handler.Tmpl = template.New(page.Name)

		handler.wvb_template_exec(wvb_config.Prefix)

		handlers[i] = handler

		mux.Handle(page.Path, handler)
		log.Print("Registered handler for " + page.Path)
	}

	//TODO go routine for re-executing templates.

	return mux
}
