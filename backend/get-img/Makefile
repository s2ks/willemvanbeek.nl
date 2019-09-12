SRC=wvb-get-img.c\
    sql-util.c

OBJS=$(SRC:.c=.o)

../../$(GET_IMG): $(OBJS)
	$(CC) $(LDFLAGS) -o $@ $(OBJS)

include $(OBJS:.o=.d)

%.d: %.c
	$(CC) $(CFLAGS) -M $< > $@

.PHONY: clean

clean:
	-rm $(OBJS) $(OBJS:%.o=%.d)
