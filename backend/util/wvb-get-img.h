#pragma once
#include <sqlite3.h>

#define sql_error() do { \
	if(sqlerr != SQLITE_OK) { \
		PRINTERR("%s\n", sqlite3_errstr(sqlerr)); \
		exit(1); \
	} \
} while(0)

#define COMMAND 0
#define DB_PATH 1
#define IMG_ID 2

#define IMG_TABLE "beelden"
#define SELECT_VAR "id"
#define SELECT_QUERY "SELECT * FROM " IMG_TABLE " WHERE id = ?;"
