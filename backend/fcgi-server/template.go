package main

import (
	"bytes"
	"html/template"
	"time"
)

//NOTE: unused
type Template interface {
	Exec(filepath string, data interface{}) (content string, err error)
	DoExec() bool
}

type FileTemplate struct {
	Name              string //template name
	File              string //file to use
	ContentQueryParam string
	Content           string //content to display
}

type PageTemplate struct {
	Prefix       string
	Template     *template.Template
	LastExec     time.Time
	ExecInterval time.Duration
	LastError    error
}

func NewFileTemplate(data *TemplateJson) (ft *FileTemplate) {
	ft = new(FileTemplate)

	ft.New(data)

	return
}

func NewPageTemplate() (pt *PageTemplate) {
	pt = new(PageTemplate)

	pt.ExecInterval = Settings.ExecInterval
	pt.Template = template.New("")

	return
}

func (pt *PageTemplate) Exec(filepath string, data interface{}, files []FileTemplate) (content string, err error) {
	var buf bytes.Buffer

	pt.Prefix = filepath
	pt.LastExec = time.Now()

	for _, file := range files {
		err = file.Parse(pt.Prefix, pt.Template)

		if err != nil {
			pt.LastError = err
			return
		}
	}

	for _, file := range files {
		err = file.Exec(&buf, data, pt.Template)
		if err != nil {
			pt.LastError = err
			return
		}
	}

	content = buf.String()

	return
}

func (pt *PageTemplate) DoExec() bool {
	if pt.LastExec.IsZero() {
		return true
	}

	if time.Now().Sub(pt.LastExec).Seconds() < pt.ExecInterval.Seconds() {
		return false
	} else {
		return true
	}
}

func (pt *PageTemplate) RegisterForExec(prefix string, data interface{}, files []FileTemplate, c chan *PageContent) {
	go func() {
		var content string
		var err error

		for {
			if !pt.DoExec() {
				time.Sleep(1 * time.Second)
				continue
			}

			content, err = pt.Exec(prefix, data, files)

			c <- &PageContent{content, err}
		}
	}()
}

func (ft *FileTemplate) New(data *TemplateJson) {

	ft.Name = data.Name
	ft.File = data.File
	ft.Content = data.Content
}

func (ft *FileTemplate) Exec(buf *bytes.Buffer, data interface{}, tmpl *template.Template) error {
	var err error

	err = tmpl.ExecuteTemplate(buf, ft.Name, data)

	return err
}

func (ft *FileTemplate) Parse(filepath string, tmpl *template.Template) error {
	var err error

	_, err = tmpl.ParseFiles(filepath + ft.File)

	return err
}
