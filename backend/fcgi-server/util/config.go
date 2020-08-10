package util

import (
	"encoding/xml"
	"fmt"
	"os"
	"log"
	"strings"
)

type ConfigTemplate struct {
	Outfile string `xml:"outfile,attr"`
	Files   []struct {
		Id   string `xml:"id,attr"`
		//Name string `xml:",innerxml"`
		Path string `xml:",innerxml"`
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

func (conf *Config) GetPageFor(path string) (*ConfigPage, error) {
	for _, page := range conf.Page {
		if page.Path == path {
			return &page, nil
		}
	}

	return nil, fmt.Errorf("No page found for path " + path)
}


func (conf *Config) Get(path string) {
	var buf []byte

	fi, err := os.Stat(path)

	if err != nil {
		log.Fatal(err)
	}

	buf = make([]byte, fi.Size())

	f, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Read(buf)

	err = xml.Unmarshal(buf, conf)

	if err != nil {
		log.Fatal(err)
	}
}

func (conf *Config) Parse(subs map[string]string) {
	/* Variable substitution */
	for key, val := range subs {
		for _, page := range conf.Page {
			page.Name = strings.ReplaceAll(page.Name, key, val)
			page.Path = strings.ReplaceAll(page.Path, key, val)
			page.Title = strings.ReplaceAll(page.Title, key, val)

			page.Template.Outfile = strings.ReplaceAll(page.Template.Outfile, key, val)

			for _, file := range page.Template.Files {
				file.Id = strings.ReplaceAll(file.Id, key, val)
				file.Path = strings.ReplaceAll(file.Path, key, val)
			}
		}
	}
}
