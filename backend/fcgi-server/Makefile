GC = go build

#fcgi server for www subdomain
WWW_SRC=www.go\
	config.go
WWW_BIN=www.fcgi

#fcgi server for admin subdomain
ADMIN_SRC=admin.go\
	  config.go
ADMIN_BIN=admin.fcgi

all: $(WWW_SRC) $(ADMIN_SRC)
	$(GC) -o ../../$(WWW_BIN) $(WWW_SRC)
	$(GC) -o ../../$(ADMIN_BIN) $(ADMIN_SRC)
