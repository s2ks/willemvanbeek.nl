package main

type WvbConfig struct {
	/*
	List of paths to handle and return a page for
	*/
	Paths []string

	/*
	template names to execute per page
	*/
	Templates [][]string

	/*
	files that define templates
	*/
	Files [][]string

	/*
	page titles
	*/
	Titles []string

	/*
	page content
	*/
	Content [][][]string
}
