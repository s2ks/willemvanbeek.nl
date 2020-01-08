package main

import (
	"database/sql"
	"log"
	"net/http"
)

const (
	dbName      = "gallery.db"
	selectField = "src"
	tableName   = "beelden"
	tableField  = "materiaal"
)

type PageGallery struct {
	Page
	Template PageTemplate

	db   *sql.DB
	stmt *sql.Stmt
	data *GalleryData
}

type GalleryData struct {
	Path  string
	Title string
	Name  string

	SrcPaths []string

	prefix string
}

func (p *PageGallery) New(page *PageJson) Handler {
	if p != nil {
		return p
	}

	return &PageGallery{
		*(NewPage(page)),
		PageTemplate{},
		nil,
		nil,
		nil,
	}
}

func (p *PageGallery) Setup(prefix string) error {
	var data GalleryData
	var err error

	data.Path = p.Path
	data.Title = p.Title
	data.Name = p.Name

	data.prefix = prefix

	p.data = &data

	p.db, err = ActiveDatabase()

	if err != nil {
		return err
	}

	p.stmt, err = p.db.Prepare("SELECT " + selectField + " FROM " + tableName + " WHERE " + tableField + "=?")

	if err != nil {
		log.Print(err)
	}

	p.Template = *(NewPageTemplate())

	return nil
}

func (p *PageGallery) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	params := q[tableField]

	p.data.SrcPaths = make([]string, 0)

	/* loop through params */
	for i := 0; i < len(params); i++ {
		rows, err := p.stmt.Query(params[i])

		if err != nil {
			log.Print(err)
			continue //DEBUG
		}

		defer rows.Close()

		for rows.Next() {
			var src string
			if err := rows.Scan(&src); err != nil {
				log.Print(err)
				continue //DEBUG
			} else {
				p.data.SrcPaths = append(p.data.SrcPaths, src)
			}
		}

	}

	if len(p.data.SrcPaths) > 0 {
		content, err := p.Template.Exec(p.data.prefix, p.data, p.Files)
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		} else {
			p.SendContent(content, err)
		}
	} else {
		log.Print("No source paths")
		http.NotFound(w, r)
	}

	p.Serve(w, r)
}
