package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type FcgiConfig struct {
	RootJson
}

func GetFcgiConfigFromProg(prog string, path string) *FcgiConfig {
	var stdout, stderr io.ReadCloser
	var err error
	var cmd *exec.Cmd

	var bytes []byte

	var fcgi_config *FcgiConfig

	fcgi_config = new(FcgiConfig)

	cmd = exec.Command(prog, path)

	stdout, err = cmd.StdoutPipe()
	log.Print(err)

	stderr, err = cmd.StderrPipe()
	log.Print(err)

	err = cmd.Start()
	if err != nil {
		log.Print(err)
		return nil
	}

	bytes, err = ioutil.ReadAll(stderr)
	if err != nil {
		log.Print(err)
		return nil
	}

	log.Print(string(bytes)) //write data received over stderr to stderr

	bytes, err = ioutil.ReadAll(stdout)
	if err != nil {
		log.Print(err)
		return nil
	}

	err = json.Unmarshal(bytes, &fcgi_config) //unmarshal data received over stdout
	if err != nil {
		log.Print(err)
		return nil
	}

	err = cmd.Wait()
	if err != nil {
		log.Print(err)
		return nil
	}

	return fcgi_config
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
