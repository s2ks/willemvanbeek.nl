package server

import (
	"encoding/xml"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"time"
	"fmt"

	//"willemvanbeek.nl/backend/server/config"
)

type FcgiServer struct {
	ServeMux *http.ServeMux

	Address  string
	Port     string
	Protocol string

	ExecInterval time.Duration

	HandlerTree *HandlerNode
}

func (s *FcgiServer) Register(path string, i IHandler) *Handler {
	mux := s.ServeMux

	handler := NewHandler(path, i)
	s.HandlerTree.Insert(handler)

	err := i.Setup(path)

	if err != nil {
		log.Fatal(err)
	}

	mux.Handle(path, handler)
	log.Print("Registered handler for " + path)

	return handler
}


func New(address string, port string, protocol string) (*FcgiServer, error) {
	var s *FcgiServer

	s = new(FcgiServer)

	s.ServeMux = http.NewServeMux()


	s.Address = address
	s.Port = port
	s.Protocol = protocol

	s.HandlerTree = NewHandlerNode(nil)

	return s, nil
}

func (s *FcgiServer) Serve() error {
	listener, err := net.Listen(s.Protocol, fmt.Sprintf("%s:%s", s.Address, s.Port))

	defer listener.Close()

	if err != nil {
		return err
	}

	handler := s.ServeMux

	return fcgi.Serve(listener, handler)

}
