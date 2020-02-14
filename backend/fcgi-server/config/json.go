package config

import (
	"encoding/json"
	"fmt"
	"os"
)

/* Defines structures to unmarshal from json file */

type NetJson struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

type SystemJson struct {
	SrvPath      string `json:"srvPath"`
	ExecInterval string `json:"templateExecutionInterval"`
}

type TemplateJson struct {
	Name              string `json:"name"`
	File              string `json:"file"`
	ContentQueryParam string `json:"contentQueryParam,omitempty"`
	Content           string `json:"content,omitempty"`
}

type PageJson struct {
	Path     string         `json:"path"`
	Title    string         `json:"title"`
	Name     string         `json:"name"`
	Display  bool           `json:"display"`
	Type     string         `json:"type,omitempty"`
	Params   []string       `json:"params,omitempty"`
	Template []TemplateJson `json:"template"`
}

type JsonConf struct {
	Net    NetJson    `json:"net"`
	System SystemJson `json:"system"`
	Page   []PageJson `json:"page"`
}

func GetJsonConf(path string) (*JsonConf, error) {
	var config *JsonConf

	var buf []byte
	var err error
	var configFile string

	configFile = path

	config = new(JsonConf)

	if configFile == "" {
		return nil, fmt.Errorf("No config file provided")
	}

	fi, err := os.Stat(configFile)

	if err != nil {
		return nil, err
	}

	/* allocate buffer equal to config file size */
	buf = make([]byte, fi.Size())

	f, err := os.Open(configFile)

	if err != nil {
		return nil, err
	}

	/* read contents of file into buf */
	_, err = f.Read(buf)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
