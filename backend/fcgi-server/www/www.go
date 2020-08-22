package main

/* www subdomain for willemvankeek.nl */

import (
	"flag"
	"os"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/s2ks/fcgiserver"
	"github.com/s2ks/fcgiserver/config"
	"github.com/s2ks/fcgiserver/logger"
)

func main() {
	var confpath = flag.String("config", "", "Path to the configuration file")
	var dbpath = flag.String("database", "", "Path to the database")
	var debug = flag.Bool("debug", false, "Enable debug logging")

	flag.Parse()

	if *debug {
		logger.LogLevel(logger.LogLevelDebug)
	}

	serverconf, err := config.GetServerConfFromXmlFile(*confpath)

	if err != nil {
		logger.Debug("In GetServerConfFromXmlFile:")
		logger.Fatal(err)
	}

	s, err := fcgiserver.New(serverconf.Net.Address, serverconf.Net.Port, serverconf.Net.Protocol)

	if err != nil {
		logger.Fatal(err)
	}

	raw, err := config.GetUserXmlFromFile(*confpath)

	if err != nil {
		logger.Debug("In GetUserXmlFromFile")
		logger.Fatal(err)
	}

	myconf, err := GetMyConfFromXml(raw)

	if err != nil {
		logger.Debug("In GetMyConfFromXml")
		logger.Fatal(err)
	}

	if _, err = os.Stat(*dbpath); err != nil {
		logger.Fatal(err)
	}

	db, err := sql.Open("sqlite3", *dbpath)

	if err != nil {
		logger.Fatal(err)
	}

	s.Register("/", &GenericPageHandler{config: myconf})
	s.Register("/contact/", &GenericPageHandler{config: myconf})
	s.Register("/beelden/steen/", &GalleryPageHandler{config: myconf, db: db, Material: "steen"})
	s.Register("/beelden/hout/", &GalleryPageHandler{config: myconf, db: db, Material: "hout"})
	s.Register("/beelden/metaal/", &GalleryPageHandler{config: myconf, db: db, Material: "metaal"})

	logger.Fatal(s.Serve())
}
