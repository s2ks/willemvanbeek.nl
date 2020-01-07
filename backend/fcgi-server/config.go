package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type FcgiConfig struct {
	RootJson
}

func GetFcgiConfig() (*FcgiConfig, error) {
	var config_file string
	var fcgiConfig *FcgiConfig
	//var data map[string]interface{}

	var buf []byte
	var err error

	fcgiConfig = new(FcgiConfig)

	/* config file will be [executable name].json */
	config_file = Settings.ConfigPath + filepath.Base(os.Args[0]) + ".json"

	fi, err := os.Stat(config_file)

	if err != nil {
		return nil, err
	}

	/* allocate buffer equal to config file size */
	buf = make([]byte, fi.Size())

	/* open file, or create if it doesn't exist */
	f, err := os.OpenFile(config_file, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}

	_, err = f.Read(buf)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, fcgiConfig)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return fcgiConfig, nil
}
