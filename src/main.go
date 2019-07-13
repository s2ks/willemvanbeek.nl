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

func intro_handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/introductie" {
		log.Print("404 " + r.URL.Path)
		http.NotFound(w, r)
		return
	}

	t, err := template.ParseFiles(
		prefix + "opt1/intro.html",
		prefix + "opt1/header.html",
		prefix + "opt1/navbar.html",
		prefix + "opt1/footer.html",
	)

	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Title string
	} {
		Title: "Carousel",
	}

	err = t.Execute(w, data)

	if err != nil {
		log.Fatal(err)
	}
}

func opt1_handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/opt1" {
		log.Print("404 " + r.URL.Path)
		http.NotFound(w, r)
		return
	}

	t, err := template.ParseFiles(
		prefix + "opt1/index.html",
		prefix + "opt1/header.html",
		prefix + "opt1/navbar.html",
		prefix + "opt1/footer.html",
	)

	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Title string
		ImageSrc []string
		ImageTitle []string
		ImageDescription []string
	} {
		Title: "Willem van Beek",
		ImageSrc: []string {
		},
		ImageTitle: []string {

		},
		ImageDescription: []string {

		},
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
		prefix + "footer.html",
	)

	if err != nil {
		log.Fatal(err)
	}
	//TODO load this data from a database
	data := struct {
		Title string
	} {
		Title: "Willem van Beek",
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
	http.HandleFunc("/opt1", opt1_handler)
	http.HandleFunc("/introductie", intro_handler)

	fcgi.Serve(listener, nil)
}
