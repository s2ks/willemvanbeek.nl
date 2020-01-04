package main

//Page types
const (
	PageTypeGeneric = "GENERIC"
	PageTypeGallery = "GALLERY"
	PageTypeAdmin	= "ADMIN"
)

/* UnMarshaled from json */
//TODO
/*type GalleryData struct {

}*/

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
