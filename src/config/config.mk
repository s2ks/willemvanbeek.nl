CC=gcc

VPATH += src/config

SRC=main.c\
    config.c

OBJS=$(SRC:.o:.c)

LDFLAGS += -lconfig
CFLAGS =

EXE = wvb.config

wvb.config: $(OBJS)
	$(CC) $(LDFLAGS) -o $(EXE) $(OBJS)

%.o : %.c
	$(CC) $(CFLAGS) -c $<
