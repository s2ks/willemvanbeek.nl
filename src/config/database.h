#pragma once
#include <stdio.h>
#include <stdarg.h>

#include <sqlite3.h>

#include "config-util.h"
#include "config.h"

#define WVB_CATEGORY_STEEN 	(1 << 0)
#define WVB_CATEGORY_HOUT	(1 << 1)
#define WVB_CATEGORY_METAAL	(1 << 2)

#define MAX_SQL_LEN 512

#define LOG_UNWRITTEN(_count) do { \
	if(_count >= MAX_SQL_LEN) \
	PRINTERR("Failed to write %d %s due to size limits\n", (_count - MAX_SQL_LEN) + 1 ,\
			(_count - MAX_SQL_LEN) + 1 == 1 ? "byte" : "bytes"); \
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

#define LAST_SQL_ERROR(db) \
	PRINTERR("%s\n", sqlite3_errmsg(db));

#define LOG_SQL_ERROR(_db, _status) do { \
	if(_status != SQLITE_OK) { \
		PRINTERR("Last sqlite api call returned %d and failed with error code %d: %s\n", _status,\
				sqlite3_errcode(_db), sqlite3_errmsg(_db)); \
	} else { \
		PRINTINFO("Last sqlite api call returned %d and completed successfully with error code %d: %s\n", \
				_status, sqlite3_errcode(_db), sqlite3_errmsg(_db)); \
	} \
} while(0)

#define IF_SQL_ERROR(_db, _status) \
	LOG_SQL_ERROR(_db, _status); \
	if(_status != SQLITE_OK)


#define CREATE_BEELDEN \
	"CREATE TABLE beelden(id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, material TEXT, "\
	"description TEXT, img_path TEXT, category INTEGER);"

#define CREATE_TABLE_FROM_BEELDEN_FMT \
	"CREATE TABLE %s AS SELECT * FROM beelden where category & %d != 0;"

#define CREATE_UPDATE_TRIGGER_FMT \
	"CREATE TRIGGER %s AFTER UPDATE ON beelden "\
	"BEGIN "\
		"DELETE FROM %s WHERE id = OLD.id;"\
		"INSERT INTO %s (id, name, description, material, img_path, category) "\
		"SELECT * FROM beelden WHERE category & %d != 0;"\
	"END;"

#define CREATE_INSERT_TRIGGER_FMT \
	"CREATE TRIGGER %s AFTER INSERT ON beelden "\
	"BEGIN "\
		"INSERT INTO %s (id, name, description, material, img_path, category) "\
		"SELECT * FROM beelden WHERE category & %d != 0;"\
	"END;"

#define CREATE_DELETE_TRIGGER_FMT \
	"CREATE TRIGGER %s DELETE ON beelden "\
	"BEGIN "\
		"DELETE FROM %s WHERE id = OLD.id;"\
	"END;"


static inline void CREATE_UPDATE_TRIGGER(sqlite3 *db, char *buf, const char *trigger_name, const char *table_name,
		int category)
{
	char *err = NULL;
	int status;

	status = snprintf(buf, MAX_SQL_LEN, CREATE_UPDATE_TRIGGER_FMT, trigger_name, table_name, table_name, category);
	LOG_UNWRITTEN(status);

	status = sqlite3_exec(db, buf, NULL, NULL, &err);
	LOG_EXEC_ERROR(status, buf, err);

	sqlite3_free(err);
}

static inline void CREATE_INSERT_TRIGGER(sqlite3 *db, char* buf, const char *trigger_name, const char *table_name,
		int category)
{
	char *err = NULL;
	int status;

	status = snprintf(buf, MAX_SQL_LEN, CREATE_INSERT_TRIGGER_FMT, trigger_name, table_name, category);
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


void wvb_database_fill_content(WVB_TEMPLATE *wvb_tmpl, sqlite3 *db);
void wvb_database_manage(sqlite3 *db);
