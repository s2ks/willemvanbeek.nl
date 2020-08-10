package main
/* www subdomain for willemvankeek.nl */
import (
	"database/sql"
	"flag"
	"log"
	"path"
	"os"
	"fmt"
	"time"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"willemvanbeek.nl/backend/server"
	"willemvanbeek.nl/backend/logger"
	"willemvanbeek.nl/backend/util"
	"willemvanbeek.nl/backend/server/config"
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
	g_db *sql.DB
	g_config *Config
)

/* Helper */
func ServeContent(w http.ResponseWriter, r *http.Request, h *server.Handler) {
	if err := h.GetErr(); err != nil {
		server.InternalServerError(w)
		logger.Error(fmt.Sprintf("Error while serving \"%s\" - %s", r.URL.Path, err))
		server.LogRequest(r)
		return
	}

	c, err := h.Content()

	if err != nil {
		server.InternalServerError(w)
		logger.Error(fmt.Sprintf("Error while fetching content for \"%s\" - %s", r.URL.Path, err))
		server.LogRequest(r)
		return
	}

	_, err = w.Write(c)

	if err != nil {
		server.InternalServerError(w)
		logger.Error(fmt.Sprintf("Error while writing response for \"%s\" - %s", r.URL.Path, err))
		server.LogRequest(r)
		return
	}
}

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

/* Execute page template */
/* TODO remove FcgiServer param */
func (p *GenericPage) Execute(s *server.FcgiServer) ([]byte, error) {
	var buf *util.Buffer
	var files []string

	buf = new(util.Buffer)

	files = make([]string, 0)

	for _, file := range p.page.Template.Files {
		files = append(files, fmt.Sprintf("%s/%s", s.TemplatePath, file.Name))
	}

	tmpl, err := template.ParseFiles(files...)
	logger.Verbose(fmt.Sprintf("%s", tmpl.DefinedTemplates()))

	if err != nil {
		return nil, err
	}

	for _, file := range p.page.Template.Files {
		err := tmpl.ExecuteTemplate(buf, file.Id, p)

		if err != nil {
			return nil, err
		}
	}

	/* Cache the result */
	if p.page.Template.Outfile != "" {
		outfile := fmt.Sprintf("%s/%s", s.Webroot, p.page.Template.Outfile)
		_, err = util.WriteToFile(outfile, buf.Bytes())

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to write to file %s -- %s", outfile, err))
		}
	}

	return buf.Bytes(), nil
}

func (p *GenericPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var buf []byte
	var size uint64
	var written uint64
	var out *util.Buffer

	out = new(util.Buffer)

	if DoServe() == false {
		//TODO 404
	}

	for _, file := range p.page.Template.Files {
		fi, err := os.Stat(file.Path)

		if err != nil {
			server.InternalServerError(w)
			//TODO log

			return
		}

		size += fi.Size()

	}

	buf = make(byte[], size)

	for _, file := range p.page.Template.Files {
		f, err := os.Open(file.Path)

		if err != nil {
			server.InternalServerError(w)
			//TODO log

			return
		}

		w, err := f.Read(buf[written:])
		written += w

		if err != nil {
			server.InternalServerError(w)
			//TODO log

			return
		}
	}

	tmpl, err := template.New(p.page.Name).Parse(string(buf[0:written]))

	if err != nil {
		server.InternalServerError(w)

		//TODO log

		return
	}

	err = tmpl.Execute(out, p)

	if err != nil {
		server.InternalServerError(w)

		//TODO log
		return
	}

	w.Write(out.Bytes())
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
/* TODO remove FcgiServer param */
func (g *GalleryPage) Execute(s *server.FcgiServer) ([]byte, error) {
	var buf *util.Buffer
	var files []string

	buf = new(util.Buffer)

	g.Thumbs = make([]string, 0)
	g.SrcPaths = make([]string, 0)

	base := path.Base(g.page.Path)

	g.Material = base

	rows, err := g.stmt.Query(path.Base(base))

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
	logger.Verbose(fmt.Sprintf("Defined templates: %s", tmpl.DefinedTemplates()))

	if err != nil {
		return nil, err
	}

	for _, file := range g.page.Template.Files {
		err := tmpl.ExecuteTemplate(buf, file.Id, g)

		if err != nil {
			return nil, err
		}
	}

	/* Cache the result */
	if g.page.Template.Outfile != "" {
		outfile := fmt.Sprintf("%s/%s", s.Webroot, g.page.Template.Outfile)
		_, err = util.WriteToFile(outfile, buf.Bytes())

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to write to file %s -- %s", outfile, err))
		}
	}

	return buf.Bytes(), nil
}

func (g *GalleryPage) ServeHTTP(w http.ResponseWriter, r *http.Request, h *server.Handler) {
	var buf []byte


	ServeContent(w, r, h)
}

func main() {
	var confpath = flag.String("config", "", "Path to the configuration file")
	var dbpath = flag.String("db", "", "Path to the database")
	var debug = flag.Bool("debug", false, "Enable debug logging")
	var conf *util.Config

	flag.Parse()

	conf = new(Config)

	serverConf, err := server.GetConfig(confpath)

	if err != nil {
		log.Fatal(err)
	}

	//s, err := server.New(*confpath, conf)
	s, err := server.New(serverConf.Net.Address, serverConf.Net.Port, serverConf.Net.Protocol)

	if err != nil {
		log.Fatal(err)
	}

	g_conf = conf

	if *debug {
		logger.LogLevel(logger.LogLevelDebug)
	}

	_, err = os.Stat(*dbpath)

	if err != nil {
		log.Fatal(err)
	}

	g_db, err = sql.Open("sqlite3", *dbpath);

	defer g_db.Close()

	if err != nil {
		log.Fatal(err)
	}

	s.Register("/", &GenericPage{})
	s.Register("/contact", &GenericPage{})
	s.Register("/beelden/steen", &GalleryPage{})
	s.Register("/beelden/hout", &GalleryPage{})
	s.Register("/beelden/metaal", &GalleryPage{})

	logger.Fatal(s.Serve())
