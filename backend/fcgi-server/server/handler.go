package server

import (
	"fmt"
	"time"
	"net/http"
)

const (
	sourceChan = 0
	relayChan  = 1

	source = 0
	relay = 1
)

type IHandler interface {
	Setup(string) error
	Execute() ([]byte, error)
	ServeHTTP(http.ResponseWriter, *http.Request, *Handler)
}

type Handler struct {
	Path     string

	channel []chan []byte
	content bool
	cache   []byte
	err     error
	ihandler IHandler
}

func (h *Handler) contentServer() {
	for {
		select {
		case c := <-h.channel[source]:
			if len(c) > 0 {
				h.cache = c
				h.content = true
			}
		default:
		}

		h.channel[relay] <- h.cache
	}
}

func NewHandler(path string, i IHandler) *Handler {
	h := new(Handler)

	h.Path = path
	h.channel = make([]chan []byte, 2)
	h.content = false

	/* Unbuffered */
	h.channel[source] = make(chan []byte, 0)
	h.channel[relay] = make(chan []byte, 0)

	h.ihandler = i

	go h.contentServer()

	return h
}

/* TODO variable timeout */
func (h *Handler) Content() ([]byte, error) {
	select {
	case c := <-h.channel[relay]:
		if h.content == false {
			return nil, fmt.Errorf("No content available for %s", h.Path)
		} else {
			return c, nil
		}
	case <-time.After(1 * time.Second):
		return nil, fmt.Errorf("Content server timeout")
	}
}

/* mutex */
var errlock = make(chan bool, 1)

func (h *Handler) GetErr() error {
	errlock <- true
	err := h.err
	<-errlock

	return err
}

func (h *Handler) SetErr(err error) {
	errlock <- true
	h.err = err
	<-errlock
}
