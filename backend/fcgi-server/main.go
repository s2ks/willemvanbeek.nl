package main

import (
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"time"
)

const(
	//ident = "willemvanbeek.nl"
	dbpath_env = "FCGI_DATABASE"
	ident_env = "FCGI_IDENT"
)


func main() {
	var err error
	var fcgiConfig *FcgiConfig
	var handler *http.ServeMux

	Settings.ConfigPath, err = os.UserConfigDir()

	if err != nil {
		log.Fatal(err)
	}

	if ident := os.Getenv(ident_env); ident == "" {
		log.Fatal("Please set " + ident_env + " environment variable")
	} else {
		Settings.ConfigPath += "/" + ident + "/"
	}


	if dbpath := os.Getenv(dbpath_env); dbpath == "" {
		Settings.DbPath = Settings.ConfigPath
	} else {
		Settings.DbPath = dbpath
	}

	/* Create config folder/file */
	if err = ConfigInit(Settings.ConfigPath); err != nil {
		log.Fatal(err)
	}

	fcgiConfig, err = GetFcgiConfig()

	if err != nil {
		log.Fatal(err)
	}

	//Interval at which templates should be re-executed
	Settings.ExecInterval, err = time.ParseDuration(fcgiConfig.System.ExecInterval)

	if err != nil {
		log.Fatal(err)
	}

	db, err := DatabaseInit(Settings.DbPath)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	wvbdb, err := ActiveDatabaseFile()

	if err != nil {
		log.Fatal(err)
	}

	err = os.Setenv("WVB_DATABASE", wvbdb)

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
