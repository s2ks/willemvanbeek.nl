package main
/* www subdomain for willemvankeek.nl */
import (
	"database/sql"
	"flag"
	"log"
	"os"
	"html/template"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"willemvanbeek.nl/backend/server"
	"willemvanbeek.nl/backend/logger"
	"willemvanbeek.nl/backend/util"
	"willemvanbeek.nl/backend/server/config"
)

type GenericPage struct {
	Title string
	Name  string

	conf *XmlConfig
	page *XmlPage
}

type GalleryPage struct {
	GenericPage
	Thumbs   []string
	SrcPaths []string

	Material string

	conf *XmlConfig
	db *sql.DB
	stmt *sql.Stmt
}

func (page *XmlPage) DoServe(r *http.Request) bool {
	p1 := strings.ToUpper(r.URL.Path)
	p2 := strings.ToUpper(page.Path)

	p1 = strings.TrimSpace(p1)
	p2 = strings.TrimSpace(p2)

	p1 = strings.TrimRight(p1, "/\\")
	p2 = strings.TrimRight(p2, "/\\")

	if p1 == p2 {
		return true
	} else {
		return false
	}
}

func (p *GenericPage) Setup(path string) error {
	page, err := p.conf.GetPageFor(path)

	if err != nil {
		return err
	}

	p.Title = page.Title
	p.Name = page.Name
	p.page = page

	return nil
}

func (p *GenericPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var out *util.Buffer

	if p.page.DoServe(r) == false {
		server.NotFound(w, r)
		server.LogRequest(r)
		return
	}

	out = new(util.Buffer)

	files := make([]string, len(p.page.Template.Files))

	for i, file := range p.page.Template.Files {
		files[i] = file.Path
	}

	buf, err := util.ReadFromFiles(files...)

	if err != nil {
		server.InternalServerError(w, r, err)
		server.LogRequest(r)
		return
	}

	tmpl, err := template.New(p.page.Name).Parse(string(buf))

	if err != nil {
		server.InternalServerError(w, r, err)
		server.LogRequest(r)
		return
	}

	err = tmpl.Execute(out, p)

	if err != nil {
		server.InternalServerError(w, r, err)
		server.LogRequest(r)
		return
	}

	w.Write(out.Bytes())
}

/* perform page setup */
/* TODO close g.db on exit */
func (g *GalleryPage) Setup(path string) error {
	page, err := g.conf.GetPageFor(path)

	if err != nil {
		return err
	}

	g.Title = page.Title
	g.Name = page.Name
	g.page = page

	/* Check if the file exists */
	_, err = os.Stat(g.page.DB.Path)

	if err != nil {
		return err
	}

	g.db, err = sql.Open("sqlite3", g.page.DB.Path)

	if err != nil {
		return err
	}

	g.stmt, err = g.db.Prepare(g.page.DB.Query)

	return err
}

func (g *GalleryPage) Scan(rows *sql.Rows) error {
	for rows.Next() {
		var src string
		var thumb string

		err := rows.Scan(&src, &thumb)

		if err != nil {
			return err
		}

		g.SrcPaths = append(g.SrcPaths, src)
		g.Thumbs = append(g.Thumbs, thumb)
	}

	err := rows.Err()

	if err != nil {
		return err
	}

	return nil
}

func (g *GalleryPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var out *util.Buffer

	if g.page.DoServe(r) == false {
		server.NotFound(w, r)
		server.LogRequest(r)
		return
	}

	out = new(util.Buffer)

	files := make([]string, len(g.page.Template.Files))

	for i, file := range g.page.Template.Files {
		files[i] = file.Path
	}

	buf, err := util.ReadFromFiles(files...)

	if err != nil {
		server.InternalServerError(w, r, err)
		server.LogRequest(r)
		return
	}

	tmpl, err := template.New(g.page.Name).Parse(string(buf))

	if err != nil {
		server.InternalServerError(w, r, err)
		server.LogRequest(r)
		return
	}

	rows, err := g.stmt.Query()

	defer rows.Close()

	if err != nil {
		server.InternalServerError(w, r, err)
		server.LogRequest(r)
		return
	}

	err = g.Scan(rows)

	if err != nil {
		server.InternalServerError(w, r, err)
		server.LogRequest(r)
		return
	}

	err = tmpl.Execute(out, g)

	if err != nil {
		server.InternalServerError(w, r, err)
		server.LogRequest(r)
		return
	}

	w.Write(out.Bytes())
}

func main() {
	var confpath = flag.String("config", "", "Path to the configuration file")
	var debug = flag.Bool("debug", false, "Enable debug logging")

	flag.Parse()

	if *debug {
		logger.LogLevel(logger.LogLevelDebug)
	}

	serverConf, err := config.GetServerConf(*confpath)

	if err != nil {
		log.Fatal(err)
	}

	s, err := server.New(serverConf.Net.Address, serverConf.Net.Port, serverConf.Net.Protocol)

	if err != nil {
		log.Fatal(err)
	}

	conf, err := GetMyConf(*confpath)

	s.Register("/", &GenericPage{ conf: conf })
	s.Register("/contact", &GenericPage{ conf: conf } )
	s.Register("/beelden/steen", &GalleryPage{ conf: conf, Material: "steen" })
	s.Register("/beelden/hout", &GalleryPage{ conf: conf, Material: "hout" })
	s.Register("/beelden/metaal", &GalleryPage{ conf: conf, Material: "metaal" })

	log.Fatal(s.Serve())
}
