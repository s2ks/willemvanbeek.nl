ROOT=www

HTTP=$(ROOT)/*
HTTP_DEST=/srv/http/

JS=node_modules/bootstrap/dist/js/* \
   node_modules/jquery/dist/*


all:$(BIN)
	$(MAKE) -C backend

.PHONY: clean dist install

dist:
	sass --update scss/custom.scss:www/css/bootstrap.css
	-cp -r $(JS) $(ROOT)/js/

clean:
	$(MAKE) -C backend clean

install:
	-cp -r $(HTTP) $(HTTP_DEST)

