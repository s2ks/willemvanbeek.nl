package util

import (
	"testing"
	"fmt"
	"encoding/xml"
	"os"
)

type VarXml struct {
	Items []struct {
		Name string `xml:"name,attr"`
		Value string `xml:",innerxml"`
	} `xml:"item"`
}

type NetXml struct {
	Address  string `xml:"address"`
	Port     string `xml:"port"`
	Protocol string `xml:"protocol"`
}

type SystemXml struct {
	Webroot  string `xml:"root"`
	Template struct {
		Path         string `xml:"path"`
		ExecInterval string `xml:"execinterval"`
	} `xml:"templates"`
}

type ServerConf struct {
	XMLName xml.Name  `xml:"server"`
	Var VarXml `xml:"vars"`
	Net     NetXml    `xml:"net"`
	System  SystemXml `xml:"system"`
}

func (conf *ServerConf) VarMap() map[string]string {
	vars := make(map[string]string)

	for _, item := range conf.Var.Items {
		vars[item.Name] = item.Value
	}

	return vars
}

func TestSub5(t *testing.T) {
	var config *ServerConf

	var buf []byte
	var err error
	var configFile string


	configFile = "./testconf.xml"

	config = new(ServerConf)

	if configFile == "" {
		t.Errorf("No config file provided")
	}

	fi, err := os.Stat(configFile)

	if err != nil {
		t.Errorf("%v", err)
	}

	buf = make([]byte, fi.Size())

	f, err := os.Open(configFile)

	if err != nil {
		t.Errorf("%v", err)
	}

	_, err = f.Read(buf)

	if err != nil {
		t.Errorf("%v", err)
	}

	/* Unmarshal for configuration variables first */
	err = xml.Unmarshal(buf, config)

	if err != nil {
		t.Errorf("%v", err)
	}

	/* Variable substitution */
	buf, err = ByteSubstituteMap(buf, config.VarMap(), "%")

	fmt.Println(string(buf))

	if err != nil {
		t.Errorf("%v", err)
	}

	err = xml.Unmarshal(buf, config)

	fmt.Println(config)

	if err != nil {
		t.Errorf("%v", err)
	}
}
