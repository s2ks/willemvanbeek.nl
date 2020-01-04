GC = gccgo

SRC=main.go \
    config.go \
    data.go \
    handler.go \
    settings.go \
    util/db.go \
    http/generic-handler.go \
    http/gallery-handler.go

$(WVB_BACKEND): $(SRC)
	$(GC) -o ../../$@ $(SRC)
