#pragma once
#include <stdio.h>
#include <stdarg.h>

#include <sqlite3.h>

#include "config-util.h"

#define WVB_CATEGORY_STEEN 	(1 << 0)
#define WVB_CATEGORY_HOUT	(1 << 1)
#define WVB_CATEGORY_METAAL	(1 << 2)

#define MAX_SQL_LEN 512

#define LOG_UNWRITTEN(_count) do { \
	if(_count > 0) \
	PRINTERR("Failed to write %d %s due to size limits\n", _count, _count == 1 ? "byte" : "bytes"); \
} while(0)
#define LOG_EXEC_ERROR(_status, _query, _errmsg) do { \
	if(_status != SQLITE_OK) { \
		PRINTWARN("Failed to execute query \"%s\" with error code %d: %s\n", \
				_query, _status, _errmsg != NULL ? _errmsg : "UNKNOWN ERROR"); \
	} \
} while(0)

#define LOG_ALLOCATION_ERROR(_dest, _bytes) do { \
	if(_dest == NULL) { \
		PRINTERR("Failed to allocate %d %s\n", (int) (_bytes), _bytes == 1 ? "byte" : "bytes"); \
		exit(-1); \
	} \
} while(0)

#define CREATE_BEELDEN "CREATE TABLE "\
	"beelden(id PRIMARY KEY, name TEXT, material TEXT, img_path TEXT, category INT);"

#define CREATE_TABLE_FROM_BEELDEN_FMT "CREATE TABLE %s AS SELECT * FROM beelden where category & %d = 1;"

#define CREATE_UPDATE_TRIGGER_FMT "CREATE TRIGGER %s AFTER UPDATE ON beelden "\
	"BEGIN "\
		"DROP TABLE IF EXISTS %s;"\
	"END;"

#define CREATE_INSERT_TRIGGER_FMT "CREATE TRIGGER %s AFTER INSERT ON beelden "\
	"BEGIN "\
		"DROP TABLE IF EXISTS %s;"\
	"END;"

#define CREATE_DELETE_TRIGGER_FMT "CREATE TRIGGER %s AFTER DELETE ON beelden "\
	"BEGIN "\
		"DROP TABLE IF EXISTS %s;"\
	"END;"


static inline void CREATE_UPDATE_TRIGGER(sqlite3 *db, char *buf, const char *trigger_name, const char *table_name)
{
	char *err = NULL;
	int status;

	status = snprintf(buf, MAX_SQL_LEN, CREATE_UPDATE_TRIGGER_FMT, trigger_name, table_name);
	LOG_UNWRITTEN(status);

	status = sqlite3_exec(db, buf, NULL, NULL, &err);
	LOG_EXEC_ERROR(status, buf, err);

	sqlite3_free(err);
}

static inline void CREATE_INSERT_TRIGGER(sqlite3 *db, char* buf, const char *trigger_name, const char *table_name)
{
	char *err = NULL;
	int status;

	status = snprintf(buf, MAX_SQL_LEN, CREATE_INSERT_TRIGGER_FMT, trigger_name, table_name);
	LOG_UNWRITTEN(status);

	status = sqlite3_exec(db, buf, NULL, NULL, &err);
	LOG_EXEC_ERROR(status, buf, err);

	sqlite3_free(err);
}

static inline void CREATE_DELETE_TRIGGER(sqlite3 *db, char* buf, const char *trigger_name, const char *table_name)
{
	char *err = NULL;
	int status;

	status = snprintf(buf, MAX_SQL_LEN, CREATE_DELETE_TRIGGER_FMT, trigger_name, table_name);
	LOG_UNWRITTEN(status);

	status = sqlite3_exec(db, buf, NULL, NULL, &err);
	LOG_EXEC_ERROR(status, buf, err);

	sqlite3_free(err);
}
