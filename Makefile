GC=gccgo

VPATH += src

SRC=main.go
OBJS=$(SRC:.go=.o)

HTTP=http/*
HTTP_DEST=/srv/http/

EXEC=wvb.backend

all: $(OBJS)
	$(GC) -o $(EXEC) $(OBJS)

%.o: %.go
	$(GC) -c $<

.PHONY: clean

clean:
	-rm $(OBJS)

install:
	-cp -r $(HTTP) $(HTTP_DEST)

