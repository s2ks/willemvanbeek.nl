package main

/* Defines structures to (un)marshal from json file */

type NetJson struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

type SystemJson struct {
	SrvPath      string `json:"srvPath"`
	ExecInterval string `json:"templateExecutionInterval"`
}

type TemplateJson struct {
	Name              string `json:"name"`
	File              string `json:"file"`
	ContentQueryParam string `json:"contentQueryParam,omitempty"`
	Content           string `json:"content,omitempty"`
}

type PageJson struct {
	Path     string         `json:"path"`
	Title    string         `json:"title"`
	Name     string         `json:"name"`
	Display  bool           `json:"display"`
	Type     string         `json:"type,omitempty"`
	Params   []string       `json:"params,omitempty"`
	Template []TemplateJson `json:"template"`
}

type RootJson struct {
	Net    NetJson    `json:"net"`
	System SystemJson `json:"system"`
	Page   []PageJson `json:"page"`
}
