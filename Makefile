ROOT=www

CC=gcc

HTTP=$(ROOT)/*
HTTP_DEST=/srv/http/

JS=node_modules/bootstrap/dist/js/* \
   node_modules/jquery/dist/* \
   node_modules/colcade/colcade.js

LDFLAGS += -lsqlite3 -ljson-c -lconfig
CFLAGS += `pkg-config --cflags json-c`
CFLAGS += -Wall -Wextra -Wpedantic
CFLAGS += -g
CFLAGS += -DDEBUG
#CFLAGS += -DVERBOSE

WVB_BACKEND 	= wvb.backend
WVB_CONFIG 	= wvb.config
GET_IMG		= get-img
ADD_IMG		= add-img

BIN = $(WVB_BACKEND) 	\
      $(WVB_CONFIG) 	\
      $(GET_IMG)	\
      $(ADD_IMG)

export LDFLAGS
export CFLAGS
export CC = gcc
export BIN
export WVB_BACKEND
export WVB_CONFIG
export GET_IMG
export ADD_IMG


.PHONY: clean dist install all

all:
	-$(MAKE) -C backend

dist:
	sass --update scss/custom.scss:www/css/bootstrap.css
	-cp -r $(JS) $(ROOT)/js/

clean:
	-$(MAKE) -C backend clean
	-rm $(BIN)

install:
	-cp -r $(HTTP) $(HTTP_DEST)

