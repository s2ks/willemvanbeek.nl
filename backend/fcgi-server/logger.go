package main

import "log"

func wvb_log(err error) bool {
	if err != nil {
		log.Print(err)
		return true
	}

	return false
}
