package main

import (
	"encoding/xml"
	"fmt"

	"willemvanbeek.nl/backend/util"
	"willemvanbeek.nl/backend/server/config"
)

type XmlTemplate struct {
	Outfile string `xml:"outfile,attr"`
	Files   []struct {
		Path string `xml:",innerxml"`
	} `xml:"file"`
}

type XmlDB struct {
	Path string `xml:"path,attr"`
	Query string `xml:"query,attr"`
}

type XmlPage struct {
	Name  string `xml:"name"`
	Path  string `xml:"path,attr"`
	Title string `xml:"title"`

	DB XmlDB `xml:"db"`
	Template XmlTemplate `xml:"template"`
}

type XmlConfig struct {
	XMLName xml.Name     `xml:"server"`
	Page    []XmlPage `xml:"page"`
}

func (conf *XmlConfig) GetPageFor(path string) (*XmlPage, error) {
	for _, page := range conf.Page {
		if page.Path == path {
			return &page, nil
		}
	}

	return nil, fmt.Errorf("No pazge found for path " + path)
}

func GetMyConf(path string) (*XmlConfig, error) {
	var conf *XmlConfig

	buf, err := util.ReadFromFile(path)

	if err != nil {
		return nil, err
	}

	vars, err := config.GetVars(path)

	if err != nil {
		return nil, err
	}

	varmap := vars.ToMap()

	buf, err = util.ByteSubstituteMap(buf, varmap, "%")

	if err != nil {
		return nil, err
	}

	conf = new(XmlConfig)

	err = xml.Unmarshal(buf, conf)

	if err != nil {
		return nil, err
	}

	return conf, nil
}
