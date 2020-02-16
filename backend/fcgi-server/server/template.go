package server

import (
	"html/template"
	"log"
	"time"
)

func (s *FcgiServer) RegisterForExec(h *Handle, data Handler) error {
	go func() {
		var lastExec time.Time

		for {
			lastExec = time.Now()
			h.Template = template.New(h.Path)

			b, err := data.Execute(h, s)

			if err != nil {
				log.Print(err)
				h.SetErr(err)
			} else {
				h.channel[sourceChan] <- b
				h.SetErr(nil)
			}

			slp := s.ExecInterval - lastExec.Sub(time.Now())
			time.Sleep(slp)
		}
	}()

	return nil
}
