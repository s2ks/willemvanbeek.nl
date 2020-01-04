package main

type ITemplate interface {
	Exec(filepath string) error
	LastError() error
	Ok() bool
}
