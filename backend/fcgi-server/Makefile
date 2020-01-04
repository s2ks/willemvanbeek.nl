GC = gccgo

SRC=main.go \
    config.go \
    data.go \
    handler.go \
    settings.go \
    logger.go \
    page.go \
    template.go\
    http/page-generic.go \
    http/page-gallery.go

$(WVB_BACKEND): $(SRC)
	$(GC) -o ../../$@ $(SRC)
