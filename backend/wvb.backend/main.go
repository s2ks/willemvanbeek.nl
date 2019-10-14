package main

import (
	"net"
	"log"
	"net/http/fcgi"
	"net/http"
	"time"
	"html/template"
	"strings"
)

var Settings struct {
	ExecInterval time.Duration
}

func NewHandler(fcgi_config *FcgiConfig) *http.ServeMux {
	var mux *http.ServeMux

	mux = http.NewServeMux()

	for _, page := range fcgi_config.Page {
		page.Type = strings.ToUpper(page.Type)

		if page.Type == PageTypeGeneric {
			handler := new(GenericHandler)
			handler.GT = new(GenericTemplate)

			handler.GT.ExecInterval = &Settings.ExecInterval

			handler.Page = &page
			handler.Path = page.Path
			handler.GT.Template = template.New(page.Name)
			handler.Display = page.Display

			handler.GenericTemplateExec(fcgi_config.Prefix)

			mux.Handle(page.Path, handler)
			log.Print("Registered handler for " + page.Path)
		}
	}

	return mux
}

//TODO get strings from json file
func main() {
	var err error
	var fcgi_config *FcgiConfig
	var handler *http.ServeMux
	var config_prog, config_path string
	var network, address string
	var exec_interval string


	//TODO get from json
	config_prog = "./wvb.config"
	config_path = "config/wvb.conf"

	network = "tcp"
	address = "localhost:9000"

	exec_interval = "10m"

	//Interval at which templates should be re-executed (refreshed)
	Settings.ExecInterval, err = time.ParseDuration(exec_interval)

	if err != nil {
		Settings.ExecInterval, _ = time.ParseDuration("10m")
	}

	listener, err := net.Listen(network, address)

	defer listener.Close()

	if err != nil {
		log.Fatal(err)
	}

	fcgi_config = GetFcgiConfig(config_prog, config_path)

	handler = NewHandler(fcgi_config)

	fcgi.Serve(listener, handler)
}
