package main

import (
	"encoding/xml"
	"fmt"
	"path"
	"net/http"
	"strings"

	"github.com/s2ks/fcgiserver/logger"
)

type XmlTemplate struct {
	Outfile string `xml:"outfile,attr"`
	Files   []struct {
		Path string `xml:",innerxml"`
	} `xml:"file"`
}

type XmlDB struct {
	Query string `xml:"query,attr"`
	Args []string `xml:"arg"`
}

type XmlPage struct {
	Name  string `xml:"name"`
	Path  string `xml:"path,attr"`
	Title string `xml:"title"`

	DB       XmlDB       `xml:"db"`
	Template XmlTemplate `xml:"template"`
}

type XmlConfig struct {
	XMLName xml.Name  `xml:"user"`
	Page    []XmlPage `xml:"page"`
}

func (page *XmlPage) DoServe(r *http.Request) bool {
	p1 := strings.ToLower(r.URL.Path)
	p2 := strings.ToLower(page.Path)

	p1 = path.Clean(p1)
	p2 = path.Clean(p2)

	logger.Debugf("%s sanitized to %s", r.URL.Path, p1)
	logger.Debugf("%s sanitized to %s", page.Path, p2)

	if p1 == p2 {
		return true
	} else {
		return false
	}
}

func (conf *XmlConfig) GetPageFor(path string) (*XmlPage, error) {
	for _, page := range conf.Page {
		if page.Path == path {
			return &page, nil
		}
	}

	return nil, fmt.Errorf("No pazge found for path " + path)
}

func GetMyConfFromXml(raw []byte) (*XmlConfig, error) {
	conf := new(XmlConfig)

	err := xml.Unmarshal(raw, conf)

	if err != nil {
		return nil, err
	}

	return conf, nil
}
