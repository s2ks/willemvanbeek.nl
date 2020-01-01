package main

import(
	"html/template"
	"time"
	"bytes"
)

type FileTemplate struct {
	Name string	//template name
	File string	//file to use
	ContentQuery string	//query used to fetch content

	Content string	//content to display
}

type PageTemplate struct {
	Prefix string
	Template *template.Template

	LastExec time.Time
	ExecInterval time.Duration

	LastError error
}

func NewFileTemplate(data *TemplateData) (ft *FileTemplate) {
	ft = new(FileTemplate)

	ft.Name = data.Name
	ft.File = data.File
	ft.ContentQuery = data.ContentQuery
	ft.Content = data.Content

	return
}

func NewPageTemplate() (pt *PageTemplate) {
	pt = new(PageTemplate)

	return
}

func (pt *PageTemplate) Exec(filepath string, data interface {}, files []FileTemplate) (err error) {
	var buf bytes.Buffer

	pt.Prefix = filepath
	pt.LastExec = time.Now()

	for _, file := range files {
		err = file.Parse(pt.Prefix)

		if err != nil {
			pt.LastError = err
			log.Print(err)
			return
		}
	}

	for _, file := range files {
		err = file.Exec(&buf, data, pt.Template) //TODO data

		if err != nil {
			pt.LastError = err
			log.Print(err)
			return
		}
	}

	pt.Content = buf.String()
}

func (pt *PageTemplate) DoExec() bool {
	if(time.Now().sub(*pt.LastExec).Seconds() < pt.ExecInterval.Seconds() {
		return false
	} else {
		return true
	}
}

func (ft *FileTemplate) Exec(bytes.Buffer buf, data interface {}, tmpl *template.Template) error {
	var err error

	_, err = tmpl.ExecuteTemplate(&buf, ft.Name, data)

	return err
}

func (ft *FileTemplate) Parse(filepath string, tmpl *template.Template) error {
	var err error

	_, err = tmpl.ParseFiles(filepath + ft.File)

	return err
}
