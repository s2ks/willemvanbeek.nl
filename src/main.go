package main

import (
	"log"
	"fmt"
	//"io"
	"net"
	"html/template"
	"net/http"
	"net/http/fcgi"
)

const prefix string = "/srv/http/template/"

func img_grid_handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/beelden/steen" {
		log.Print("404 " + r.URL.Path)
		http.NotFound(w, r)
		return
	}

	t, err := template.ParseFiles(
		prefix + "opt1/grid.html",
		prefix + "opt1/header.html",
		prefix + "opt1/navbar.html",
		prefix + "opt1/footer.html",
	)

	if err != nil {
		log.Print(err)
		http.NotFound(w, r)
		return
	}

	data := struct {
		Title string
		Content [][]string
	} {
		Title: "Willem van Beek",
		Content: [][]string {
			{
			"/static/img/20160909_125120.jpg",
			"/static/img/20160909_125120.jpg",
			"/static/img/20160909_125120.jpg",
			},
			{
			"Foto",
			"Foto",
			"Foto",
			},
			{
			"text",
			"text",
			"text",
			},
		},
	}

	err = t.Execute(w, data)

	if err != nil {
		log.Print(err)
		http.NotFound(w, r)
		return
	}
}

func intro_handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/introductie" {
		log.Print("404 " + r.URL.Path)
		http.NotFound(w, r)
		return
	}

	t, err := template.ParseFiles(
		prefix + "intro.html",
		prefix + "header.html",
		prefix + "navbar.html",
		prefix + "footer.html",
	)

	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Title string
		PageName string
	} {
		Title: "Willem van Beek",
		PageName: "Introductie",
	}

	err = t.Execute(w, data)

	if err != nil {
		log.Fatal(err)
	}
}

func root_handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		fmt.Println("404 " + r.URL.Path)
		http.NotFound(w, r)
		return
	}

	//TODO load html file path from a config file
	t, err := template.ParseFiles(
		prefix + "index.html",
		prefix + "header.html",
		prefix + "navbar.html",
		prefix + "footer.html",
	)

	if err != nil {
		log.Fatal(err)
	}
	//TODO load this data from a database
	data := struct {
		Title string
		PageName string
	} {
		Title: "Willem van Beek",
		PageName: "Home",
	}

	err = t.Execute(w, data)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:9000")

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", root_handler)
	http.HandleFunc("/introductie", intro_handler)
	http.HandleFunc("/beelden/steen", img_grid_handler)

	fcgi.Serve(listener, nil)
}
