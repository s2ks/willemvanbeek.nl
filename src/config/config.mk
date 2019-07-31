CC=gcc

VPATH += src/config

SRC=main.c\
    config.c\
    config-util.c\
    database.c

OBJS=$(SRC:.c=.o)

LDFLAGS += -lconfig
LDFLAGS += -lsqlite3
LDFLAGS += `pkg-config --libs json-c`
CFLAGS +=`pkg-config --cflags json-c`
CFLAGS += -Wall -Wextra -Wpedantic
CFLAGS += -g
CFLAGS += -DDEBUG

EXE = wvb.config

wvb.config: $(OBJS)
	$(CC) $(LDFLAGS) -o $(EXE) $(OBJS)

include $(OBJS:.o=.d)

%.d : %.c
	$(CC) $(CFLAGS) -M $< > $@
