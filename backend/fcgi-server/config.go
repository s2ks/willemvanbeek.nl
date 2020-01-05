package main

import (
	"log"
	"io/ioutil"
	"io"
	"os/exec"
	"encoding/json"
	"path/filepath"
	"os"
)

type FcgiConfig struct {
	Prefix string	//prepended to File in FileTemplate
	Database string //unused

	Page []PageData
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

//TODO implement
func GetFcgiConfig() *FcgiConfig {
	var config_file string
	var fcgi_config *FcgiConfig
	//var data map[string]interface{}

	var b []byte
	var err error

	fcgi_config = new(FcgiConfig)


	config_file = Settings.ConfigPath + filepath.Base(os.Args[0]) + ".json"

	f, err := os.OpenFile(config_file, os.O_RDWR | os.O_CREATE, 0666)

	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Read(b)

	if err != nil {
		log.Fatal(err)
	}

	log.Print(b)

	//err = json.Unmarshal()

	//fcgi_config.Prefix =

	/*
	see https://gobyexample.com/json

	search file for parameters, and fill FcgiConfig struct
	*/

	return fcgi_config
}
