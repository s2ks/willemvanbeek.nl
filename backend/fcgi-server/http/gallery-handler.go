package main

import(
	"net/http"
	"bytes"
	"time"
	"html/template"
	"log"
	"encoding/json"
)

//TODO inherit GenericHandler
type GalleryHandler struct {
	GenericHandler

	T *GalleryTemplate
}

//TODO inherit GenericTemplate
type GalleryTemplate struct {
	GenericTemplate
}

//TODO inherit GenericTemplateData
type GalleryTemplateData struct {
	GenericTemplateData

	ContentQuery []string
	ImgData []ImgData
}


//TODO DRY
func NewGalleryHandler(page *PageData, path string, prefix string, display bool, execInterval *time.Duration) (h *GalleryHandler) {
	h = new(GalleryHandler)
	h.T = new(GalleryTemplate)

	h.Page = page
	h.Path = path
	h.T.Prefix = prefix
	h.Display = display
	h.T.ExecInterval = *execInterval

	h.TemplateExec(prefix)

	return
}

//TODO implement
func (d *GalleryTemplateData) QueryImages() {
	for i, query := range d.ContentQuery {
		if query == "" {
			continue
		}

		data, err := ExecQuery(query)
		if err != nil {
			log.Print(err)
			continue
		}
		//TODO ImgData struct
		err = json.Unmarshal(data.stdout, &d.ImgData[i])
	}
}

//TODO DRY
func (h *GalleryHandler) TemplateExec(prefix string) (err error) {
	var buf bytes.Buffer

	data := GalleryTemplateData {
		h.Page.Path,
		h.Page.Title,
		h.Page.Name,
		"",
		nil,
		nil,
	}

	err = nil

	h.T.Prefix = prefix
	h.T.LastExec = time.Now()

	data.ContentQuery = make([]string, len(h.Page.Template))
	data.ImgData = make([]ImgData, len(h.Page.Template))


	for i, tmpl := range h.Page.Template {
		_, err = h.T.Template.ParseFiles(prefix + tmpl.File)

		data.ContentQuery[i] = tmpl.ContentQuery

		if err != nil {
			h.T.LastError = err
			log.Print(err)
			return
		}
	}

	data.QueryImages()

	for _, tmpl := range h.Page.Template {
		err = h.T.Template.ExecuteTemplate(&buf, tmpl.Name, data)

		if err != nil {
			h.T.LastError = err
			log.Print(err)
			return
		}
	}

	h.Content = buf.String()

	return
}

//TODO DRY
func (h *GalleryHandler) GetHandlerData() HandlerData {
	return HandlerData {
		h.Path,
		h.Content,
		h.T.Prefix,
		h.Display,
	}
}

//TODO DRY
func (h *GalleryHandler) LastError() error {
	return h.T.LastError
}

//TODO DRY
func (h *GalleryHandler) CheckTime() bool {
	return CheckTime(&h.T.LastExec, &h.T.ExecInterval)
}

//TODO DRY
func (h *GalleryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	HandleServeHTTP(w, r, h)
}


