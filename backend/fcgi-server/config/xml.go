package config

import (
	"encoding/xml"
	"fmt"
	"os"
)

/* Defines structures to unmarshal from xml config file */

type TemplateSystemXml struct {
	Path         string `xml:"path"`
	ExecInterval string `xml:"execinterval"`
}

type PageTemplateFileXml struct {
	Id   string `xml:"id,attr"`
	Post string `xml:"post,attr"`
	File string `xml:",innerxml"`
}

type PageTemplateXml struct {
	OutFile string                `xml:"outfile,attr"`
	Files   []pageTemplateFileXml `xml:"file"`
}

type NetXml struct {
	Address string `xml:"address"`
	Port    string `xml:"port"`
	Proto   string `xml:"protocol"`
}

type SystemXml struct {
	WebRoot  string            `xml:"webroot"`
	Template templateSystemXml `xml:"template"`
}

type PageXml struct {
	Name  string `xml:"name,attr"`
	Path  string `xml:"path"`
	Title string `xml:"title"`
	Type  string `xml:"type"`

	Template []PageTemplateXml `xml:"template"`
}

type XmlConf struct {
	Net    netXml    `xml:"net"`
	System systemXml `xml:"system"`
	Page   []pageXml `xml:"page"`
}

func GetXmlConf(path string) (*XmlConf, error) {
	var config *XmlConf

	var buf []byte
	var err error
	var configFile string

	configFile = path

	config = new(XmlConf)

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
