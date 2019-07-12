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


func root_handler(w http.ResponseWriter, r *http.Request) {

	prefix := "/srv/http/template/"
	if r.URL.Path != "/" {
		fmt.Println(r.URL.Path)
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
		Title: "Beelden | Willem van Beek",

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

	fcgi.Serve(listener, nil)
}
