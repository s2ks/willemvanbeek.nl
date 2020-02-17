package main

import (
	"bytes"
	"database/sql"
	"flag"
	"log"
	"path"
	"os"
	"fmt"
	"html/template"

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

func (p *GenericPage) Execute(s *server.FcgiServer) ([]byte, error) {
	var buf *bytes.Buffer
	var files []string

	buf = new(bytes.Buffer)

	files = make([]string, len(p.page.Template.Files))

	for i, file := range p.page.Template.Files {
		files[i] = fmt.Sprintf("%s/%s", s.TemplatePath, file.Name)
		log.Print(fmt.Sprintf("\t%s - %s", p.page.Path, files[i]))
	}
	tmpl, err := template.ParseFiles(files...)

	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(buf, p)

	if err != nil {
		return nil, err
	} else {
		log.Print(fmt.Sprintf("\texecuted %s - %s", p.page.Path, string(buf.Bytes())))
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
func (g *GalleryPage) Execute(s *server.FcgiServer) ([]byte, error) {
	var buf *bytes.Buffer
	var files []string

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


	files = make([]string, len(g.page.Template.Files))

	for i, file := range g.page.Template.Files {
		files[i] = fmt.Sprintf("%s/%s", s.TemplatePath, file.Name)
		log.Print(fmt.Sprintf("\t%s - %s", g.page.Path, files[i]))
	}

	tmpl, err := template.ParseFiles(files...)

	if err != nil {
		return nil, err
	}

	err = tmpl.Execute(buf, g)

	if err != nil {
		return nil, err
	} else {
		log.Print(fmt.Sprintf("\texecuted %s - %s", g.page.Path, string(buf.Bytes())))
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
