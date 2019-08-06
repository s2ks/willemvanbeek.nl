GC = gccgo

VPATH += src/

EXE = wvb.backend

GOSRC=src/main.go \
    src/config.go \
    src/logger.go \
    src/page-handler.go

wvb.backend: $(GOSRC)
	$(GC) -o $@ $(GOSRC)

