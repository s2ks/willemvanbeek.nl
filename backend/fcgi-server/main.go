package main

import (
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"time"
)

const ident = "willemvanbeek.nl"

func main() {
	var err error
	var fcgiConfig *FcgiConfig
	var handler *http.ServeMux

	Settings.ConfigPath, err = os.UserConfigDir()

	if err != nil {
		log.Fatal(err)
	}

	Settings.ConfigPath += "/" + ident + "/"

	/* Check if willemvanbeek.nl folder exists, if not create it*/
	if _, err = os.Stat(Settings.ConfigPath); os.IsNotExist(err) {
		err = os.Mkdir(Settings.ConfigPath, 0777)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	fcgiConfig, err = GetFcgiConfig()

	if err != nil {
		log.Fatal(err)
	}

	//Interval at which templates should be re-executed (refreshed)
	Settings.ExecInterval, err = time.ParseDuration(fcgiConfig.System.ExecInterval)

	if err != nil {
		log.Fatal(err)
	}

	address := fcgiConfig.Net.Address + ":" + fcgiConfig.Net.Port

	listener, err := net.Listen(fcgiConfig.Net.Protocol, address)

	defer listener.Close()

	if err != nil {
		log.Fatal(err)
	}

	handler = NewHandler(fcgiConfig)

	log.Fatal(fcgi.Serve(listener, handler))
}
