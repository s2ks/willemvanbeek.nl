GC = gccgo

BIN = wvb.backend

FILES=main.go \
    config.go \
    logger.go \
    page-handler.go

LDFLAGS += -lsqlite3 -ljson-c -lconfig
CFLAGS += `pkg-config --cflags json-c`
CFLAGS += -Wall -Wextra -Wpedantic
CFLAGS += -g
CFLAGS += -DDEBUG
#CFLAGS += -DVERBOSE

SRC=$(FILES:%=src/%)

export LDFLAGS
export CFLAGS
export CC = gcc

VPATH += src

all: $(BIN)
	$(MAKE) -C config
	$(MAKE) -C util

$(BIN): $(SRC)
	$(GC) -o $@ $(SRC)

.PHONY: clean

clean:
	$(MAKE) -C config clean
	$(MAKE) -C util clean
	-rm $(BIN)
