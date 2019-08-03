GC=gccgo

VPATH += src

GOSRC=src/main.go \
      src/config.go \
      src/logger.go
GOOBJS=wvb.backend.o

HTTP=http/*
HTTP_DEST=/srv/http/

DIST=webpack-wvb/dist/js \
     webpack-wvb/dist/css
DIST_DEST=http/res/

BIN=wvb.backend\
    wvb.config

all:$(BIN)

wvb.backend: $(GOSRC)
	$(GC) -o $@ $(GOSRC)

include src/config/config.mk


.PHONY: clean dist install

dist:
	-cd webpack && yarn run build

clean:
	-rm $(OBJS) $(GOOBJS) $(OBJS:.o=.d)

install:
	-cp -r $(DIST) $(DIST_DEST)
	-cp -r $(HTTP) $(HTTP_DEST)

