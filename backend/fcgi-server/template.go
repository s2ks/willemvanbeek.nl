package main

import(
	"html/template"
	"time"
	"bytes"
)

type Template interface {
	Exec(filepath string, data interface{}) (content string, err error)
	DoExec() bool
}

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

	ft.New(data)

	return
}

func NewPageTemplate() (pt *PageTemplate) {
	pt = new(PageTemplate)

	pt.ExecInterval = Settings.ExecInterval

	return
}


func (pt *PageTemplate) Exec(filepath string, data interface {}, files []FileTemplate) (content string, err error) {
	var buf bytes.Buffer

	pt.Prefix = filepath
	pt.LastExec = time.Now()

	for _, file := range files {
		err = file.Parse(pt.Prefix)

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

	if time.Now().sub(*pt.LastExec).Seconds() < pt.ExecInterval.Seconds() {
		return false
	} else {
		return true
	}
}

func (pt *PageTemplate) RegisterForExec(prefix string, data interface{}, c chan string) {
	go func() {
		var content string
		var err error

		for {
			if !pt.DoExec() {
				time.Sleep(1 * time.Second)
				continue
			}

			content, err = pt.Exec(prefix, data)

			if err != nil {
				log.Print(err)
				continue
			}

			c <- content
		}
	}()
}

func (ft *FileTemplate) New(data *TemplateData) {

	ft.Name = data.Name
	ft.File = data.File
	ft.ContentQuery = data.ContentQuery
	ft.Content = data.Content
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
