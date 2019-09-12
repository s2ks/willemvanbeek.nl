SRC=main.c\
    config.c\
    config-util.c\
    database.c


OBJS=$(SRC:.c=.o)

../../$(WVB_CONFIG): $(OBJS)
	$(CC) $(LDFLAGS) -o $@ $(OBJS)

include $(OBJS:.o=.d)

%.d : %.c
	$(CC) $(CFLAGS) -M $< > $@

.PHONY: clean

clean:
	-rm $(OBJS) $(OBJS:%.o=%.d)
