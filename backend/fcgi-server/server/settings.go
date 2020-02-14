package main

import (
	"time"

	"willemvanbeek.nl/backend/config"
)

var Settings struct {
	ExecInterval time.Duration
	QueryProg    string
	DbPath       string
	ConfigPath   string
	Config *config.XmlConf
}
