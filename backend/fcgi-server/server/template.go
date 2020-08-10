package server

import (
	"log"
	"time"
)

func RegisterForExec(h *Handler, interval time.Duration) error {
	go func() {
		var lastExec time.Time

		for {
			lastExec = time.Now()

			b, err := h.ihandler.Execute()

			if err != nil {
				log.Print(err)
				h.SetErr(err)
			} else {
				h.channel[sourceChan] <- b
				h.SetErr(nil)
			}

			slp := interval - lastExec.Sub(time.Now())
			time.Sleep(slp)
		}
	}()

	return nil
}
