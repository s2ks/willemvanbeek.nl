package main

import (
	"log"
	"io/ioutil"
	"io"
	"os/exec"
	"encoding/json"
)

type FcgiConfig struct {
	Prefix string	//prepended to File in WvbTemplate
	Database string //unused

	Page []PageData
}

func GetFcgiConfig(prog string, path string) *FcgiConfig {
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
