package config

import (
	"encoding/xml"
	"fmt"
	"os"
)

/* Defines structures to unmarshal from xml config file */

type NetXml struct {
	Address  string `xml:"address"`
	Port     string `xml:"port"`
	Protocol string `xml:"protocol"`
}

type SystemXml struct {
	Webroot  string `xml:"webroot"`
	Template struct {
		Path         string `xml:"path"`
		ExecInterval string `xml:"execinterval"`
	} `xml:"template"`
}

type ServerConf struct {
	XMLName xml.Name  `xml:"server"`
	Net     NetXml    `xml:"net"`
	System  SystemXml `xml:"system"`
}

func GetServerConf(path string) (*ServerConf, error) {
	var config *ServerConf

	var buf []byte
	var err error
	var configFile string

	configFile = path

	config = new(ServerConf)

	if configFile == "" {
		return nil, fmt.Errorf("No config file provided")
	}

	fi, err := os.Stat(configFile)

	if err != nil {
		return nil, err
	}

	buf = make([]byte, fi.Size())

	f, err := os.Open(configFile)

	if err != nil {
		return nil, err
	}

	_, err = f.Read(buf)

	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(buf, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
