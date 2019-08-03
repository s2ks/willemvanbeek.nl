package main

import "log"

func wvb_log(err error) int {
	if err != nil {
		log.Print(err)
		return 1
	}

	return 0
}

func printerr(bytes []byte) {
	log.Print(string(bytes))
}
