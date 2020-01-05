package main

import (
	"net"
	"log"
	"net/http/fcgi"
	"net/http"
	"time"
	"os"
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
		log.Fatal(err)
	}
	Settings.ConfigPath, err = os.UserConfigDir()

	if err != nil {
		log.Fatal(err)
	}

	Settings.ConfigPath += "/willemvanbeek.nl/"

	/* Check if willemvanbeek.nl folder exists, if not create it*/
	if _, err = os.Lstat(Settings.ConfigPath); os.IsNotExist(err) {
		err = os.Mkdir(Settings.ConfigPath, 0777)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen(network, address)

	defer listener.Close()

	if err != nil {
		log.Fatal(err)
	}

	fcgi_config = GetFcgiConfigFromProg(config_prog, config_path)

	GetFcgiConfig()

	return //DEBUG

	handler = NewHandler(fcgi_config)

	log.Fatal(fcgi.Serve(listener, handler))
}
