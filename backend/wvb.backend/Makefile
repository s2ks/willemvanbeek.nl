GC = gccgo

SRC=main.go \
    config.go \
    logger.go \
    page-handler.go

$(WVB_BACKEND): $(SRC)
	$(GC) -o ../../$@ $(SRC)
