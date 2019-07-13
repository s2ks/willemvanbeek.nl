GC=gccgo

VPATH += src

SRC=main.go
OBJS=$(SRC:.go=.o)

HTTP=http/*
HTTP_DEST=/srv/http/

DIST=webpack-wvb/dist/js \
     webpack-wvb/dist/css
DIST_DEST=http/res/

EXEC=wvb.backend

all: $(OBJS)
	$(GC) -o $(EXEC) $(OBJS)

%.o: %.go
	$(GC) -c $<

.PHONY: clean dist install

dist:
	-cd webpack-wvb && yarn run build

clean:
	-rm $(OBJS)

install:
	-cp -r $(DIST) $(DIST_DEST)
	-cp -r $(HTTP) $(HTTP_DEST)

