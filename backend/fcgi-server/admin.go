package main

import (
	"net/http"

	"willemvanbeek.nl/backend/logger"
	"willemvanbeek.nl/backend/server"
)

type AdminPage struct {
	page *ConfigPage
}

type LoginPage struct {
	ActiveLogins int

	page *ConfigPage
}

func (a *AdminPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if Auth == false {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	server.ServeHTTP(w, r)
}

func main() {


	s, err := server.New()

	if err != nil {
		logger.Fatal(err)
	}

	s.Register("/login", &LoginPage{})

	s.Register("/", &AdminPage{})
	s.Register("/add", &AdminPage{})
	s.Register("/edit", &AdminPage{})
	s.Register("/delete", &AdminPage{})

	logger.Fatal(s.Serve())
}
