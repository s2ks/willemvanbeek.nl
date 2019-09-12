package main

import (
	"net"
	"log"
	"net/http/fcgi"
)

func main() {
	var wvb_config *WvbConfig

	listener, err := net.Listen("tcp", "localhost:9000")

	defer listener.Close()

	if err != nil {
		log.Fatal(err)
	}

	wvb_config = wvb_config_fetch("backend/wvb.config", "config/wvb.conf")

	handler := wvb_handler_init(wvb_config)

	fcgi.Serve(listener, handler)
}
