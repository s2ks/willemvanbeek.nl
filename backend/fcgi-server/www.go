package main

import (
	"bytes"
	"database/sql"
	"flag"
	"log"
	"path"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"willemvanbeek.nl/backend/server"
)

type GenericPage struct {
	Title string
	Name  string

	page *ConfigPage
}

type GalleryPage struct {
	GenericPage
	Thumbs   []string
	SrcPaths []string

	stmt *sql.Stmt
}

var (
	db *sql.DB
	configData *Config
)

func (p *GenericPage) Setup(path string) error {
	page, err := GetPageFor(path)

	if err != nil {
		return err
	}

	p.Title = page.Title
	p.Name = page.Name
	p.page = page

	return nil
}

func (p *GenericPage) Execute(h *server.Handle, s *server.FcgiServer) ([]byte, error) {
	var buf *bytes.Buffer

	buf = new(bytes.Buffer)

	for _, file := range p.page.Template.Files {
		_, err := h.Template.ParseFiles(s.TemplatePath + file.Name)

		if err != nil {
			return nil, err
		}
	}
	err := h.Template.Execute(buf, p)

	if err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

/* perform page setup */
func (g *GalleryPage) Setup(path string) error {
	page, err := GetPageFor(path)

	if err != nil {
		return err
	}

	g.Title = page.Title
	g.Name = page.Name
	g.page = page

	g.stmt, err = db.Prepare("SELECT src, thumb FROM gallery WHERE type = ?")

	return err
}

/* Execute page template */
func (g *GalleryPage) Execute(h *server.Handle, s *server.FcgiServer) ([]byte, error) {
	var buf *bytes.Buffer

	buf = new(bytes.Buffer)

	g.Thumbs = make([]string, 0)
	g.SrcPaths = make([]string, 0)

	rows, err := g.stmt.Query(path.Base(g.page.Path))

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var src string
		var thumb string

		if err := rows.Scan(&src, &thumb); err != nil {
			return nil, err
		}

		g.SrcPaths = append(g.SrcPaths, src)
		g.Thumbs = append(g.Thumbs, thumb)
	}

	for _, file := range g.page.Template.Files {
		_, err := h.Template.ParseFiles(s.TemplatePath + file.Name)

		if err != nil {
			return nil, err
		}
	}

	err = h.Template.Execute(buf, g)

	if err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func main() {
	var confpath = flag.String("config", "", "Path to the configuration file")
	var dbpath = flag.String("db", "", "Path to the database")
	var conf *Config

	flag.Parse()

	conf = new(Config)

	s, err := server.New(*confpath, conf)

	if err != nil {
		log.Fatal(err)
	}
	configData = conf

	_, err = os.Stat(*dbpath)

	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("sqlite3", *dbpath);

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	s.Register("/", &GenericPage{})
	s.Register("/contact", &GenericPage{})

	s.Register("/beelden/steen", &GalleryPage{})
	s.Register("/beelden/hout", &GalleryPage{})
	s.Register("/beelden/metaal", &GalleryPage{})

	log.Fatal(s.Serve())
}
