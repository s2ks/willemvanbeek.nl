package main

//Page types
const (
	PageTypeGeneric = "GENERIC"
	PageTypeGallery = "GALLERY"
	PageTypeAdmin	= "ADMIN"
)

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


//summary of all data that can be used in a template
type Data struct {
	Img []ImageData
	Page PageData
}
