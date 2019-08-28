#include <limits.h>
#include <stdlib.h>
#include <errno.h>
#include <string.h>
#include <stdio.h>

#include <sqlite3.h>
#include <json.h>

#include "printerr.h"

#define COMMAND 0
#define DB_PATH 1
#define IMG_ID 2

int main(int argc, char **argv)
{
	char *dbpath, *endptr;
	int id;

	if(argc < 3) {
		printf("Usage: %s <path> <id>\n", argv ? argv[COMMAND] : "<command>");
		exit(1);
	}

	dbpath = argv[DB_PATH];

	errno = 0;
	id = strtol(argv[IMG_ID], &endptr, 10);

	if(errno) {
		PRINTERR("Error converting id to a numerical value: %s\n", strerror(errno));
		exit(errno);
	}

	if(endptr == argv[IMG_ID]) {
		PRINTERR("<id>: No digits were found\n");
		exit(1);
	}

	if(*endptr)
		PRINTWARN("<id>: Characters after digits '%d': '%s'\n", id, endptr);

	return 0;
}
