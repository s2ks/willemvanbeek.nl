package server

import (
	"log"
	"time"
)

func (s *FcgiServer) RegisterForExec(h *Handle, data Handler) error {
	go func() {
		var lastExec time.Time

		for {
			lastExec = time.Now()

			b, err := data.Execute(s)

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
