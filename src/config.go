package main

import (
	"log"
	"io/ioutil"
	"io"
	"os/exec"
	"encoding/json"
)

type WvbTemplate struct {
	Name string	//template name
	File string	//file to use
	ContentQuery string	//query used to fetch content

	Content [][]string	//content to display
}

type WvbPage struct {
	Path string	//url to handle
	Title string	//page title
	Name string	//page name

	Display bool

	Template []WvbTemplate
}

type WvbConfig struct {
	Prefix string	//prepended to File in WvbTemplate
	Database string //unused

	Page []WvbPage
}

/*
	Start config_prog with argument config_path and
	open a pipe to stdout and stderr.

	config_prog should write data in json format to stdout and any errors, warnings,
	or debug info to stderr.

	The json data received is unmarshaled to the WvbConfig struct, and returned.
*/
func wvb_config_fetch(config_prog string, config_path string) *WvbConfig {
	var stdout, stderr io.ReadCloser
	var err error
	var cmd *exec.Cmd

	var bytes []byte

	var wvb_config *WvbConfig

	wvb_config = new(WvbConfig)

	cmd = exec.Command(config_prog, config_path)

	stdout, err = cmd.StdoutPipe()
	wvb_log(err)

	stderr, err = cmd.StderrPipe()
	wvb_log(err)

	err = cmd.Start()
	if wvb_log(err) == true {
		return nil
	}

	bytes, err = ioutil.ReadAll(stderr)
	if wvb_log(err) == true {
		return nil
	}

	log.Print(string(bytes)) //write data received over stderr to stderr

	bytes, err = ioutil.ReadAll(stdout)
	if wvb_log(err) == true {
		return nil
	}

	err = json.Unmarshal(bytes, &wvb_config) //unmarshal data received over stdout
	if wvb_log(err) == true {
		return nil
	}


	err = cmd.Wait()
	if wvb_log(err) == true {
		return nil
	}

	return wvb_config
}
