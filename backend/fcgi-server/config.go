package main

import (
	"encoding/xml"
	"fmt"
)

type ConfigTemplate struct {
	Outfile string `xml:"outfile,attr"`
	Files   []struct {
		Id   string `xml:"id,attr"`
		Name string `xml:",innerxml"`
	} `xml:"file"`
}

type ConfigPage struct {
	Name  string `xml:"name"`
	Path  string `xml:"path,attr"`
	Title string `xml:"title"`

	Template ConfigTemplate `xml:"template"`
}

type Config struct {
	XMLName xml.Name     `xml:"server"`
	Page    []ConfigPage `xml:"page"`
}

func GetPageFor(path string) (*ConfigPage, error) {
	for _, page := range configData.Page {
		if page.Path == path {
			return &page, nil
		}
	}

	return nil, fmt.Errorf("No page found for path " + path)
}
