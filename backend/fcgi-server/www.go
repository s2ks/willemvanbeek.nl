package main

import (
	"net/http"
	"database/sql"
	"strings"
	"html/template"

	"willemvanbeek.nl/backend/server"
)

type GenericPage struct {
	Title string
	Name string
}

type GalleryPage struct {
	GenericPage
	Thumbs []string
	SrcPaths []string

	stmt *sql.Stmt
}

func (p *GenericPage) Execute(t *template.Template, id string, post []server.TemplatePost) ([]byte, error) {
	return server.Execute(t, id, post, p)
}

/* perform page setup */
func (g *GalleryPage) Setup() error {
	db, err := server.ActiveDB()

	if err != nil {
		return err
	}

	g.stmt, err := db.Prepare("SELECT src, thumbnail FROM beelden WHERE ?=?")

	return err
}

/* Execute page template */
func (g *GalleryPage) Execute(t *template.Template, id string, post []server.TemplatePost) ([]byte, error) {

	g.Thumbs = make([]string, 0)
	g.SrcPaths = make([]string, 0)

	for _, q := range post {
		rows, err := g.stmt.Query(q.Arg, q.Val)

		if err != nil {
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var src string
			var thumb string

			if err := rows.Scan(&src, &thumb); err != nil{
				return nil, err
			}

			g.SrcPaths = append(g.SrcPaths, src)
			g.Thumbs = append(g.Thumbs, thumb)
		}
	}

	return server.Execute(t, id, post, g)
}


func main() {} {
	server.Init()

	server.Register("/", &GenericPage{})
	server.Register("/contact", &GenericPage{})

	server.Register("/beelden/steen", &GalleryPage{});
	server.Register("/beelden/hout", &GalleryPage{});
	server.Register("/beelden/metaal", &GalleryPage{});

	server.Start()
}
