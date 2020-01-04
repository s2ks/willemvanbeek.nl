package main

import(
	"net/http"
	"bytes"
	"time"
	"html/template"
	"log"
)
type PageGeneric struct {
	Page
	Template PageTemplate

	ContentChannel chan string
}

type GenericData struct {
	Path string
	Title string
	Name string
}

func (p *PageGeneric) Setup(prefix string) {
	/* TODO
	- Initial execution of template
	- Register template for re-execution at interval
	- (Template Executor class? use go-routines)
	- Template interface
	- Populate GenericData struct with GenericPage data for template
	*/

	var data GenericData
	var err error

	p.ContentChannel = make(chan string)

	data.Path = p.Path
	data.Title = p.Title
	data.Name = p.Name

	p.Template.ExecInterval = Settings.ExecInterval
	p.Template.RegisterForExec(prefix, data, p.ContentChannel)
}


func (p *PageGeneric) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	select {
	case c := <-p.ContentChannel:
		p.SetContent(c)
	default:
	}

	p.Serve(w, r)
}
