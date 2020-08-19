package main

import (
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/s2ks/fcgiserver"
	"github.com/s2ks/fcgiserver/util"
	"github.com/s2ks/fcgiserver/logger"
)

type GenericPageHandler struct {
	Title string
	Name  string

	config *XmlConfig
	page   *XmlPage
}

type GalleryPageHandler struct {
	Title string
	Name  string

	Thumbs   []string
	SrcPaths []string

	Material string //TODO rename to Tag

	config *XmlConfig
	page   *XmlPage

	db   *sql.DB
	stmt *sql.Stmt
}

func (h *GenericPageHandler) Setup(path string) error {
	page, err := h.config.GetPageFor(path)

	if err != nil {
		return err
	}

	h.Title = page.Title
	h.Name = page.Name
	h.page = page

	return nil
}

func (h *GenericPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var out *util.Buffer

	if h.page.DoServe(r) == false {
		fcgiserver.NotFound(w, r)
		fcgiserver.LogRequest(r)
		return
	}

	logger.Debugf("Serving: %s", r.URL.Path)

	out = new(util.Buffer)

	files := make([]string, len(h.page.Template.Files))

	for i, file := range h.page.Template.Files {
		files[i] = file.Path
	}

	buf, err := util.ReadFromFiles(files...)

	if err != nil {
		fcgiserver.InternalServerError(w, r, err)
		fcgiserver.LogRequest(r)
		return
	}

	tmpl, err := template.New(h.page.Name).Parse(string(buf))

	if err != nil {
		fcgiserver.InternalServerError(w, r, err)
		fcgiserver.LogRequest(r)
		return
	}

	err = tmpl.Execute(out, h)

	if err != nil {
		fcgiserver.InternalServerError(w, r, err)
		fcgiserver.LogRequest(r)
		return
	}

	w.Write(out.Bytes())
}

/* perform page setup */
func (h *GalleryPageHandler) Setup(path string) error {
	page, err := h.config.GetPageFor(path)

	if err != nil {
		return err
	}

	h.Title = page.Title
	h.Name = page.Name
	h.page = page



	h.stmt, err = h.db.Prepare(h.page.DB.Query)

	return err
}

func (h *GalleryPageHandler) Scan(rows *sql.Rows) error {
	for rows.Next() {
		var src string
		var thumb string

		err := rows.Scan(&src, &thumb)

		if err != nil {
			return err
		}

		h.SrcPaths = append(h.SrcPaths, src)
		h.Thumbs = append(h.Thumbs, thumb)
	}

	err := rows.Err()

	if err != nil {
		return err
	}

	return nil
}

func (h *GalleryPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var out *util.Buffer

	if h.page.DoServe(r) == false {
		fcgiserver.NotFound(w, r)
		fcgiserver.LogRequest(r)
		return
	}

	logger.Debugf("Serving: %s", r.URL.Path)

	out = new(util.Buffer)

	files := make([]string, len(h.page.Template.Files))

	for i, file := range h.page.Template.Files {
		files[i] = file.Path
	}

	buf, err := util.ReadFromFiles(files...)

	if err != nil {
		fcgiserver.InternalServerError(w, r, err)
		fcgiserver.LogRequest(r)
		return
	}

	logger.Debugf("Parsing template: %s", string(buf))

	tmpl, err := template.New(h.page.Name).Parse(string(buf))

	if err != nil {
		fcgiserver.InternalServerError(w, r, err)
		fcgiserver.LogRequest(r)
		return
	}

	/* Convert string array to an interface array for Query() */
	args := make([]interface{}, len(h.page.DB.Args))

	for i, v := range h.page.DB.Args {
		args[i] = v
	}

	rows, err := h.stmt.Query(args...)

	if err != nil {
		fcgiserver.InternalServerError(w, r, err)
		fcgiserver.LogRequest(r)
		return
	}

	defer rows.Close()

	err = h.Scan(rows)

	if err != nil {
		fcgiserver.InternalServerError(w, r, err)
		fcgiserver.LogRequest(r)
		return
	}

	err = tmpl.Execute(out, h)

	if err != nil {
		fcgiserver.InternalServerError(w, r, err)
		fcgiserver.LogRequest(r)
		return
	}

	logger.Debugf("Writing content: %s", string(out.Bytes()))

	w.Write(out.Bytes())
}
