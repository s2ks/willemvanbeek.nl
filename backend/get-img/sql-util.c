#include <stdint.h>
#include <stdio.h>

#include <sqlite3.h>
#include <json.h>

json_object *get_sql_val(sqlite3_stmt *stmt, int i)
{
	json_object *json_obj = NULL;
	char *blob;

	switch(sqlite3_column_type(stmt, i)) {
		case SQLITE_INTEGER:
			json_obj = json_object_new_int64(sqlite3_column_int64(stmt, i));
			break;
		case SQLITE_FLOAT:
			json_obj = json_object_new_double(sqlite3_column_double(stmt, i));
			break;
		case SQLITE_BLOB:
			blob = malloc(64);
			if(blob == NULL) {
				json_obj = NULL;
				break;
			}
			sprintf(blob, "%p", sqlite3_column_blob(stmt, i));
			json_obj = json_object_new_string((const char *) blob);
			break;
		case SQLITE_NULL:
			json_obj = NULL;
			break;
		case SQLITE_TEXT:
			json_obj = json_object_new_string((const char *) sqlite3_column_text(stmt, i));
			break;
		default:
			json_obj = NULL;
	};

	return json_obj;
}

json_object *sql_to_json(sqlite3_stmt *stmt)
{
	int data;
	const char *column_name;
	json_object *json_obj, *val, *json_arr;

	json_arr = json_object_new_array();

	while(sqlite3_step(stmt) == SQLITE_ROW) {
		json_obj = json_object_new_object();
		data = sqlite3_data_count(stmt);

		for(int i = 0; i < data; i++) {
			column_name = sqlite3_column_origin_name(stmt, i);

			if(column_name == NULL)
				goto err;

			val = get_sql_val(stmt, i);
			json_object_object_add(json_obj, column_name, val);
		}
		json_object_array_add(json_arr, json_obj);
	}

	return json_obj;
err:
	return NULL;
}
