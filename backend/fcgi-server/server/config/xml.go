package config

import (
	"encoding/xml"
	"fmt"
	"os"

	"willemvanbeek.nl/backend/util"
)

/* Defines structures to unmarshal from xml config file */

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

type VarsXml struct {
	XMLName xml.Name `xml:"server"`
	Var VarXml `xml:"vars"`
}

func GetVars(confpath string) (*VarsXml, error) {
	var vars *VarsXml

	buf, err := util.ReadFromFile(confpath)

	if err != nil {
		return nil, err
	}

	vars = new(VarsXml)

	err = xml.Unmarshal(buf, vars)

	if err != nil {
		return nil, err
	}

	return vars, nil
}

func (vars *VarsXml) ToMap() map[string]string {
	varmap := make(map[string]string)

	for _, item := range vars.Var.Items {
		varmap[item.Name] = item.Value
	}

	return varmap
}

func (conf *ServerConf) VarMap() map[string]string {
	vars := make(map[string]string)

	for _, item := range conf.Var.Items {
		vars[item.Name] = item.Value
	}

	return vars
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

	/* Unmarshal for configuration variables first */
	err = xml.Unmarshal(buf, config)

	if err != nil {
		return nil, err
	}

	/* Variable substitution */
	buf, err = util.ByteSubstituteMap(buf, config.VarMap(), "%")

	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(buf, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
