GC = gccgo

SRC=main.go \
    config.go \
    data.go \
    data/generic-data.go \
    http/generic-handler.go

$(WVB_BACKEND): $(SRC)
	$(GC) -o ../../$@ $(SRC)
