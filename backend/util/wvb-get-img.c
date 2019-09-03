#include <limits.h>
#include <stdlib.h>
#include <errno.h>
#include <string.h>
#include <stdio.h>

#include <sqlite3.h>
#include <json.h>

#include "printerr.h"
#include "wvb-get-img.h"
#include "sql-util.h"

int main(int argc, char **argv)
{
	char *dbpath, *endptr;
	int id, sqlerr;
	sqlite3 *db;
	sqlite3_stmt *stmt;

	json_object *json;

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

	sqlerr = sqlite3_open_v2((const char *) dbpath, &db, SQLITE_OPEN_READONLY, NULL);
	if(sqlerr != SQLITE_OK || db == NULL) {
		PRINTERR("error opening %s: %s\n", dbpath, sqlite3_errstr(sqlerr));
		exit(1);
	}

	sqlerr = sqlite3_prepare_v2(db, SELECT_QUERY, -1, &stmt, NULL);
	if(sqlerr != SQLITE_OK) {
		PRINTERR("error preparing statement %s: %s error code %d\n", SELECT_QUERY, sqlite3_errstr(sqlerr), sqlerr);
		exit(1);
	}

	sqlerr = sqlite3_bind_int(stmt, 1, id);
	sql_error();

	json = sql_to_json(stmt);

	if(json == NULL) {
		PRINTERR("%s\n", sqlite3_errmsg(db));
		exit(1);
	}

	fprintf(stdout, "%s\n", json_object_to_json_string_ext(json, JSON_C_TO_STRING_NOSLASHESCAPE
#ifdef DEBUG
				| JSON_C_TO_STRING_PRETTY
				| JSON_C_TO_STRING_PRETTY_TAB
#endif
				));

	return 0;
}
