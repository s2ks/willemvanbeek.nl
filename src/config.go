package main

import (
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
	if wvb_log(err) == 1 {
		return nil
	}

	bytes, err = ioutil.ReadAll(stderr)
	wvb_log(err)

	printerr(bytes)

	bytes, err = ioutil.ReadAll(stdout)
	wvb_log(err)

	err = json.Unmarshal(bytes, &wvb_config)
	if wvb_log(err) == 1 {
		return nil
	}


	err = cmd.Wait()
	wvb_log(err)

	return wvb_config
}
