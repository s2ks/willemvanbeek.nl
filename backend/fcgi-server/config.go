package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type FcgiConfig struct {
	RootJson
}

func GetFcgiConfig() (*FcgiConfig, error) {
	var fcgiConfig *FcgiConfig
	//var data map[string]interface{}

	var buf []byte
	var err error
	var configFile string

	configFile = Settings.ConfigPath

	fcgiConfig = new(FcgiConfig)

	if configFile == "" {
		return nil, fmt.Errorf("Config has not been initialised")
	}

	fi, err := os.Stat(configFile)

	if err != nil {
		return nil, err
	}

	/* allocate buffer equal to config file size */
	buf = make([]byte, fi.Size())

	f, err := os.OpenFile(configFile, os.O_RDWR, 0666)

	if err != nil {
		return nil, err
	}

	/* read contents of file into buf */
	_, err = f.Read(buf)

	if err != nil {
		return nil, err
	}

	/* unmarshal buf into fcgiConfig */
	err = json.Unmarshal(buf, fcgiConfig)

	if err != nil {
		return nil, err
	}

	return fcgiConfig, nil
}
