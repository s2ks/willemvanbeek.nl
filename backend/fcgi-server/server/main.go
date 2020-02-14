package server

import (
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"time"
	"flag"

	"willemvanbeek.nl/backend/config"
	"willemvanbeek.nl/backend/database"
)

const (
	//ident = "willemvanbeek.nl"
	dbpath_env = "FCGI_DATABASE"
	ident_env  = "FCGI_IDENT"
	config_env = "FCGI_CONFIG"
)

func Init() {
	var err error
	var ok bool
	var conf *config.XmlConf
	var handler *http.ServeMux

	var confpath = flag.String("config", "", "Path to the configuration file")
	var dbpath = flag.String("db", "", "Path to the database")

	if confpath == "" {
		confpath, ok = os.LookupEnv(config_env)

		if ok == false {
			log.Fatal("Configuration file not found")
		}

		Settings.ConfigPath = confpath

	}

	if dbpath == "" {
		dbpath, ok = os.LookupEnv(dbpath_env)

		if ok == false {
			log.Fatal("Databse file not found")
		}

		Settings.DbPath = dbpath
	}

	conf, err = config.GetXmlConf(confpath)

	if err != nil {
		log.Fatal(err)
	}

	Settings.Config = conf

	/* Interval at which templates should be re-executed */
	Settings.ExecInterval, err = time.ParseDuration(conf.System.ExecInterval)

	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Open(dbpath)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	address := conf.Net.Address + ":" + conf.Net.Port

	listener, err := net.Listen(conf.Net.Protocol, address)

	defer listener.Close()

	if err != nil {
		log.Fatal(err)
	}

	handler = NewHandler(conf)

	log.Fatal(fcgi.Serve(listener, handler))
}

func Start() {
	
}
