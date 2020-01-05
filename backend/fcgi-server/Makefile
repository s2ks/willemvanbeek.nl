GC = go build

SRC=main.go \
    config.go \
    data.go \
    handler.go \
    settings.go \
    logger.go \
    page.go \
    template.go\
    page-generic.go \
    page-gallery.go

$(WVB_BACKEND): $(SRC)
	$(GC) -o ../../$@ $(SRC)
