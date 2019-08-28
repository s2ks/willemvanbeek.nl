SRC=main.c\
    config.c\
    config-util.c\
    database.c

BIN=../wvb.config

OBJS=$(SRC:.c=.o)

all: $(BIN)

$(BIN): $(OBJS)
	$(CC) $(LDFLAGS) -o $@ $(OBJS)

include $(OBJS:.o=.d)

%.d : %.c
	$(CC) $(CFLAGS) -M $< > $@

.PHONY: clean

clean:
	-rm $(OBJS) $(OBJS:%.o=%.d) $(BIN)
