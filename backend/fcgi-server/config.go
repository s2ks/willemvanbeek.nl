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

var configFile string

func ConfigInit(configPath string) error {
	var err error

	/* create configPath if needed */
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		if err = os.MkdirAll(configPath, 0777); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	/* configFile will be [executable name].json */
	if exec, err := os.Executable(); err != nil {
		return err
	} else {
		configFile = configPath + filepath.Base(exec) + ".json"
	}

	/* create configFile if needed */
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if _, err = os.Create(configFile); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func GetFcgiConfig() (*FcgiConfig, error) {
	var fcgiConfig *FcgiConfig
	//var data map[string]interface{}

	var buf []byte
	var err error

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
