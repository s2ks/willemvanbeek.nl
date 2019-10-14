package main

type GenericPageTemplate struct {
	Name string	//template name
	File string	//file to use
	ContentQuery string	//query used to fetch content

	Content string	//content to display
}

type GenericPage struct {
	Path string	//url to handle
	Title string	//page title
	Name string	//page name

	Display bool

	Type string

	Template []GenericPageTemplate
}
