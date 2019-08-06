ROOT=www

HTTP=$(ROOT)/*
HTTP_DEST=/srv/http/

JS=node_modules/bootstrap/dist/js/* \
   node_modules/jquery/dist/*

BIN=wvb.backend\
    wvb.config

all:$(BIN)

include src/config/config.mk
include src/backend.mk


.PHONY: clean dist install

dist:
	sass --update scss/custom.scss:www/css/bootstrap.css
	-cp -r $(JS) $(ROOT)/js/

clean:
	-rm $(OBJS) $(OBJS:.o=.d)

install:
	-cp -r $(HTTP) $(HTTP_DEST)

