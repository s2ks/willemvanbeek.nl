GC=gccgo
CC=gcc

VPATH += src
VPATH += src/config

GOSRC=src/main.go \
      src/config.go
GOOBJS=wvb.backend.o
CSRC=config.c
COBJS=$(CSRC:.c=.o)

HTTP=http/*
HTTP_DEST=/srv/http/

DIST=webpack-wvb/dist/js \
     webpack-wvb/dist/css
DIST_DEST=http/res/

EXEC=wvb.backend

all:$(GOSRC) $()
	$(GC) -o $(EXEC) $(GOSRC)
	$(CC) -lconfig -ljson-c -o test

%.o: %.c
	$(CC) -I/usr/include/json-c -c $<

.PHONY: clean dist install

dist:
	-cd webpack-wvb && yarn run build

clean:
	-rm $(OBJS)

install:
	-cp -r $(DIST) $(DIST_DEST)
	-cp -r $(HTTP) $(HTTP_DEST)

