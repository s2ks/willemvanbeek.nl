package main

import (
	"bytes"
	"html/template"
	"time"
	"fmt"
	"strings"
	"os"

	"willemvanbeek.nl/backend/config"
)

type TemplatePost struct {
	Key string
	Val string
}

func TemplateFor(path string) (*config.PageTemplateXml, error) {
	for _, page := range Settings.Config.Page {
		if page.Path == path {
			return &page.Template, nil
		}
	}

	return nil, fmt.Errorf("No template found for " + path)
}

func RegisterForExec(path string, h Handler, interval time.Duration c chan []byte) error {
	var buf *bytes.Buffer
	ct, err := TemplateFor(path)

	if err != nil {
		return err
	}

	post := make([]TemplatePost 0)
	t := template.New(path)

	for _, file := range ct.Files {
		if file.Post != "" {
			/* split key=val&key=val into key=val pairs */
			split := strings.Split(file.Post, "&")

			for _, entry := range split {
				/* split key=val into key and val */
				keyval := strings.Split(entry, "=")

				if len(split) < 2 {
					continue
				}

				post = append(post, TemplatePost {Key: keyval[0], Val: keyval[1]}
			}

			_, err := t.ParseFiles(Settings.Config.System.Template.Path + file.Name)
		}
	}

	go func() {
		for {
			buf.Reset()

			for _, file := range ct.Files {
				b, err := h.Execute(file.Name, file.Id, post, t)

				if err != nil {
					log.Print(err)
				} else {
					buf.Write(b)
				}
			}

			if ct.OutFile != "" {
				f, err := os.Create(Settings.Conf.System.Webroot + ct.OutFile)

				if err != nil {
					log.Print(err)
					goto Sleep
				}
				if _, err := f.Write(buf); err != nil {
					log.Print(err)
					goto Sleep
				}
				if err := f.Close(); err != nil {
					log.Print(err)
					goto Sleep
				}
			}


			Sleep:
			c <- buf.Bytes()

			time.Sleep(interval)
		}
	}()

	return nil
}

func Execute(t *template.Template, id string, post []TemplatePost, h Handler) ([]byte, error) {
	var buffer *bytes.Buffer

	err := t.ExecuteTemplate(buffer, id, h)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}




/*-----------------------------------------------*/
/*-----------------------------------------------*/
/*-----------------------------------------------*/
/*-----------------------------------------------*/
/*-----------------------------------------------*/
/*-----------------------------------------------*/
/*-----------------------------------------------*/
/*-----------------------------------------------*/
/*-----------------------------------------------*/



func (pt *PageTemplate) Exec(filepath string, data interface{}, files []FileTemplate) (content string, err error) {
	var buf bytes.Buffer

	pt.Prefix = filepath
	pt.LastExec = time.Now()

	pt.Template = template.New("content")

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
