#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include <sqlite3.h>

#include "config-util.h"
#include "config.h"
#include "database.h"

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
}

//fill wvb_tmpl->content with data from db
void wvb_database_fill_content(WVB_TEMPLATE *wvb_tmpl, sqlite3 *db)
{
	int status, column_count, row_count,  row, bytes;
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

	//FIXME this is a possible memory leak, maybe free() first if not null
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
