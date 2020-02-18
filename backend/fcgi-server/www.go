package main

import (
	"database/sql"
	"flag"
	"log"
	"path"
	"os"
	"fmt"
	"html/template"

	_ "github.com/mattn/go-sqlite3"

	"willemvanbeek.nl/backend/server"
	"willemvanbeek.nl/backend/logger"
	"willemvanbeek.nl/backend/util"
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

	Material string

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
	var buf *util.Buffer
	var files []string

	buf = new(util.Buffer)

	files = make([]string, 0)

	for _, file := range p.page.Template.Files {
		files = append(files, fmt.Sprintf("%s/%s", s.TemplatePath, file.Name))
	}

	tmpl, err := template.ParseFiles(files...)
	logger.Verbose(fmt.Sprintf("%s",tmpl.DefinedTemplates()))

	if err != nil {
		return nil, err
	}

	for _, file := range p.page.Template.Files {
		err := tmpl.ExecuteTemplate(buf, file.Id, p)

		if err != nil {
			return nil, err
		}
	}

	outfile := fmt.Sprintf("%s/%s", s.Webroot, p.page.Template.Outfile)

	_, err = util.WriteToFile(outfile, buf.Bytes())

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to write to file %s -- %s", outfile, err))
	}


	return buf.Bytes(), nil
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
	var buf *util.Buffer
	var files []string

	buf = new(util.Buffer)

	g.Thumbs = make([]string, 0)
	g.SrcPaths = make([]string, 0)

	base := path.Base(g.page.Path)

	g.Material = base

	rows, err := g.stmt.Query(path.Base(base)

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


	files = make([]string, 0)

	for _, file := range g.page.Template.Files {
		files = append(files, fmt.Sprintf("%s/%s", s.TemplatePath, file.Name))
	}

	tmpl, err := template.ParseFiles(files...)
	logger.Verbose(fmt.Sprintf("%s", tmpl.DefinedTemplates()))

	if err != nil {
		return nil, err
	}

	for _, file := range g.page.Template.Files {
		err := tmpl.ExecuteTemplate(buf, file.Id, g)

		if err != nil {
			return nil, err
		}
	}

	outfile := fmt.Sprintf("%s/%s", s.Webroot, g.page.Template.Outfile)

	_, err = util.WriteToFile(outfile, buf.Bytes())

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to write to file %s -- %s", outfile, err))
	}

	return buf.Bytes(), nil
}

func main() {
	var confpath = flag.String("config", "", "Path to the configuration file")
	var dbpath = flag.String("db", "", "Path to the database")
	var debug = flag.Bool("debug", false, "Enable debug logging")
	var conf *Config

	flag.Parse()

	conf = new(Config)

	s, err := server.New(*confpath, conf)

	if err != nil {
		log.Fatal(err)
	}
	configData = conf

	if *debug {
		logger.LogLevel(logger.LogLevelDebug)
	}

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

	logger.Fatal(s.Serve())
}
