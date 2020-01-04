package main

//Page types
const (
	PageTypeGeneric = "GENERIC"
	PageTypeGallery = "GALLERY"
	PageTypeAdmin	= "ADMIN"
)


/*
type PageData struct {
	Path string
	Title string
	Name string

	Action []string
	Method string

	Content [][]string
}

type ImageData struct {
	Name string
	Description string
	Source string

	Material []string
	HasMaterial func(string) bool

	Group string //reserved
}
*/


/* UnMarshaled from json */
type GalleryData struct {
	
}

type TemplateData struct {
	Name string	//template name
	File string	//file to use
	ContentQuery string	//query used to fetch content

	Content string	//content to display
}

type PageData struct {
	Path string	//url to handle //TODO rename to Url
	Title string	//page title
	Name string	//page name

	Display bool

	Type string

	Template []TemplateData
}
