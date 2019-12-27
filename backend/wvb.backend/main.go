package main

import (
	"net"
	"log"
	"net/http/fcgi"
	"net/http"
	"time"
)


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
	Settings.QueryProg = "./query-db"

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
