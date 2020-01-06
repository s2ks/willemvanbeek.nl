package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	//sqlite3 "github.com/mattn/go-sqlite3"
)

type QueryData struct {
	cmd    *exec.Cmd
	stdout []byte
	stderr []byte
}

func ExecQuery(query string) (data *QueryData, err error) {
	var stdout, stderr io.ReadCloser

	data = new(QueryData)
	err = nil

	data.cmd = exec.Command(Settings.QueryProg, query)
	data.cmd.Env = append(os.Environ(), "DATABASE="+Settings.DbPath)

	stdout, err = data.cmd.StdoutPipe()

	if err != nil {
		log.Print(err)
		return
	}

	stderr, err = data.cmd.StderrPipe()

	if err != nil {
		log.Print(err)
		return
	}

	err = data.cmd.Start()
	if err != nil {
		log.Print(err)
		return
	}

	data.stdout, err = ioutil.ReadAll(stdout)

	if err != nil {
		log.Print(err)
		return
	}

	data.stderr, err = ioutil.ReadAll(stderr)

	if err != nil {
		log.Print(err)
		return
	}

	err = data.cmd.Wait()

	if err != nil {
		log.Print(err)
		return
	}

	return
}
