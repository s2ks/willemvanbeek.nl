package server

import (
	"encoding/xml"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"time"

	"willemvanbeek.nl/backend/server/config"
)

type Handler interface {
	Setup(string) error
	Execute(*Handle, *FcgiServer) ([]byte, error)
}

type FcgiServer struct {
	ServeMux *http.ServeMux

	Address  string
	Port     string
	Protocol string

	Webroot      string
	TemplatePath string
	ExecInterval time.Duration

	Handles []*Handle
}

func (s *FcgiServer) Register(path string, data Handler) {
	mux := s.ServeMux

	handle := NewHandle(path)
	s.Handles = append(s.Handles, handle)

	err := data.Setup(path)

	if err != nil {
		log.Fatal(err)
	}

	s.RegisterForExec(handle, data)

	mux.Handle(path, handle)
	log.Print("Registered handler for " + path)
}

func ProcUserConf(path string, data interface{}) error {
	var buf []byte

	fi, err := os.Stat(path)

	if err != nil {
		return err
	}

	buf = make([]byte, fi.Size())

	f, err := os.Open(path)

	if err != nil {
		return err
	}

	_, err = f.Read(buf)

	if err != nil {
		return err
	}

	/* unmarshal buf into data */
	err = xml.Unmarshal(buf, data)

	return err
}

func New(configPath string, data interface{}) (*FcgiServer, error) {
	var s *FcgiServer

	s = new(FcgiServer)

	s.ServeMux = http.NewServeMux()

	err := ProcUserConf(configPath, data)

	if err != nil {
		return nil, err
	}

	conf, err := config.GetServerConf(configPath)

	if err != nil {
		return nil, err
	}

	s.Address = conf.Net.Address
	s.Port = conf.Net.Port
	s.Protocol = conf.Net.Protocol

	s.Webroot = conf.System.Webroot
	s.TemplatePath = conf.System.Template.Path
	s.ExecInterval, err = time.ParseDuration(conf.System.Template.ExecInterval)

	if err != nil {
		return nil, err
	}

	s.Handles = make([]*Handle, 0)

	return s, nil
}

func (s *FcgiServer) Serve() error {
	listener, err := net.Listen(s.Protocol, s.Address)

	defer listener.Close()

	if err != nil {
		return err
	}

	handler := s.ServeMux

	return fcgi.Serve(listener, handler)

}
