#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include <sqlite3.h>

#include "config-util.h"
#include "config.h"
#include "database.h"

static void wvb_create_tables(sqlite3 *db)
{
	char *query = NULL, *err = NULL;
	int status;

	query = malloc(MAX_SQL_LEN + sizeof(char));
	if(query == NULL) {
		PRINTERR("wvb_create_tables failed to allocate %ld bytes\n", MAX_SQL_LEN + sizeof(char));
		return;
	}

	//Create table "steen"
	status = snprintf(query, MAX_SQL_LEN, CREATE_TABLE_FROM_BEELDEN_FMT, "steen", WVB_CATEGORY_STEEN);
	LOG_UNWRITTEN(status);

	status = sqlite3_exec(db, query, NULL, NULL, &err);
	LOG_EXEC_ERROR(status, query, err);

	//Create table "hout"
	status = snprintf(query, MAX_SQL_LEN, CREATE_TABLE_FROM_BEELDEN_FMT, "hout", WVB_CATEGORY_HOUT);
	LOG_UNWRITTEN(status);

	status = sqlite3_exec(db, query, NULL, NULL, &err);
	LOG_EXEC_ERROR(status, query, err);

	//Create table "metaal"
	status = snprintf(query, MAX_SQL_LEN, CREATE_TABLE_FROM_BEELDEN_FMT, "metaal", WVB_CATEGORY_METAAL);
	LOG_UNWRITTEN(status);

	status = sqlite3_exec(db, query, NULL, NULL, &err);
	LOG_EXEC_ERROR(status, query, err);

	free(query);
}

static void wvb_create_triggers(sqlite3 *db)
{
	char *query = NULL;

	query = malloc(MAX_SQL_LEN + sizeof(char));
	if(query == NULL) {
		PRINTERR("Failed to allocate %ld bytes\n", MAX_SQL_LEN + sizeof(char));
		return;
	}

	/*
	 * Create triggers to drop tables if the "main" table is updated
	 * TODO find a better solution to storing beelden to be displayed on a page in a database
	 * */
	CREATE_UPDATE_TRIGGER(db, query, "steen_update_trigger", "steen", WVB_CATEGORY_STEEN);
	CREATE_UPDATE_TRIGGER(db, query, "hout_update_trigger", "hout", WVB_CATEGORY_HOUT);
	CREATE_UPDATE_TRIGGER(db, query, "metaal_update_trigger", "metaal", WVB_CATEGORY_METAAL);

	CREATE_INSERT_TRIGGER(db, query, "steen_insert_trigger", "steen", WVB_CATEGORY_STEEN);
	CREATE_INSERT_TRIGGER(db, query, "hout_insert_trigger", "hout", WVB_CATEGORY_HOUT);
	CREATE_INSERT_TRIGGER(db, query, "metaal_insert_trigger", "metaal", WVB_CATEGORY_METAAL);

	CREATE_DELETE_TRIGGER(db, query, "steen_delete_trigger", "steen");
	CREATE_DELETE_TRIGGER(db, query, "hout_delete_trigger", "hout");
	CREATE_DELETE_TRIGGER(db, query, "metaal_delete_trigger", "metaal");

	free(query);
}

void wvb_database_manage(sqlite3 *db)
{
	char *err = NULL;
	int status;

	status = sqlite3_exec(db, CREATE_BEELDEN, NULL, NULL, &err);
	if(status != SQLITE_OK) {
		if(err != NULL)
			PRINTWARN("CREATE_BEELDEN failed with error code %d: %s\n", status, err);
		else
			PRINTERR("CREATE_BEELDEN: Unknown error %d.\n", status);
	}

	sqlite3_free(err);

	wvb_create_tables(db);
	wvb_create_triggers(db);
}

//fill wvb_tmpl->content with data from db
void wvb_database_fill_content(WVB_TEMPLATE *wvb_tmpl, sqlite3 *db)
{
	int status, column_count, row_count, row, bytes;
	sqlite3_stmt *stmt = NULL;

	const unsigned char *text;
	char *buf;

	status = sqlite3_prepare_v2(db, wvb_tmpl->content_query, -1, &stmt, NULL);

	if(stmt == NULL) {
		PRINTERR("Failed to prepare query: \"%s\"\n", wvb_tmpl->content_query);
		LOG_SQL_ERROR(db, status);
		return;
	}
	column_count = sqlite3_column_count(stmt);

	if(column_count == 0) {
		LAST_SQL_ERROR(db);
		goto err;
	}

	//FIXME could leak memory
	wvb_tmpl->content = NULL;

	row_count = row = 0;
	while((status = sqlite3_step(stmt)) != SQLITE_DONE) {
		if(status != SQLITE_ROW)
			LOG_SQL_ERROR(db, status);

		wvb_tmpl->content = realloc(wvb_tmpl->content, ++row_count * sizeof(void *));
		LOG_ALLOCATION_ERROR(wvb_tmpl->content, row_count * sizeof(void *));

		wvb_tmpl->content[row] = malloc((column_count + 1) * sizeof(void *));
		LOG_ALLOCATION_ERROR(wvb_tmpl->content[row], (column_count + 1) * sizeof(void *));

		for(int i = 0; i < column_count; i++) {
			text = sqlite3_column_text(stmt, i);
			bytes = sqlite3_column_bytes(stmt, i);

			if(bytes == 0) {
				LAST_SQL_ERROR(db);
				wvb_tmpl->content[row][i] = "";
				continue;
			}

			bytes += sizeof(char);
			buf = malloc(bytes);
			LOG_ALLOCATION_ERROR(buf, bytes);

			wvb_tmpl->content[row][i] = buf;

			//write text to buf
			snprintf(buf, bytes, "%s", text);
		}
		wvb_tmpl->content[row++][column_count] = NULL;
	}

	//null-terminate row
	wvb_tmpl->content = realloc(wvb_tmpl->content, ++row_count * sizeof(void *));
	wvb_tmpl->content[row] = NULL;
err:
	sqlite3_finalize(stmt);
}
