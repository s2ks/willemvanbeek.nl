#include <stdio.h>
#include <string.h>

#include <json.h>
#include <sqlite3.h>

#include "config.h"
#include "config-util.h"
#include "database.h"

//FIXME unused
#define WVB_INDEXOF_NOTFOUND -1

static sqlite3 *db;

static inline json_object *wvb_new_string_or_null(const char *str)
{
	if(str)
		return json_object_new_string(str);
	else
		return NULL;
}

/*
 * Get the index of the page where path == wvb_config->path
 */
//FIXME unused
static int wvb_indexof_path(WVB_CONFIG *wvb_config, const char *path)
{
	for(int i = 0; i < wvb_config->page_count; i++)
		if(strcmp(path, wvb_config->page[i].path) == 0)
			return i;

	return WVB_INDEXOF_NOTFOUND;
}

//FIXME content is now an array of an array of strings
//	each array is NULL-terminated and content_length is unused
static void wvb_free_content(WVB_TEMPLATE *wvb_tmpl)
{
	if(wvb_tmpl->content == NULL)
		return;

	for(int i = 0; wvb_tmpl->content[i] != NULL; i++) {
		for(int x = 0; wvb_tmpl->content[i][x] != NULL; x++)
			free(wvb_tmpl->content[i][x]);
		free(wvb_tmpl->content[i]);
	}
	free(wvb_tmpl->content);
}

/*
 * fill WVB_TEMPLATE's content array using content_query
 */
int wvb_query_content(WVB_TEMPLATE *wvb_tmpl)
{
	if(wvb_tmpl->content_query == NULL)
		return 0;

	if(wvb_tmpl->content)
		wvb_free_content(wvb_tmpl);

	//allocate and fill WVB_TEMPLATE.content array
	wvb_database_fill_content(wvb_tmpl, db);

	return 1;
}
/*
 * convert a WVB_TEMPLATE struct to a json object
 */
json_object *wvb_template_to_json_object(WVB_TEMPLATE *wvb_tmpl)
{
	json_object *obj, *array, *array2, *val;

	obj = json_object_new_object();

	//template identifier name
	val = wvb_new_string_or_null(wvb_tmpl->name);
	json_object_object_add(obj, "name", val);

	//template html file
	val = wvb_new_string_or_null(wvb_tmpl->file);
	json_object_object_add(obj, "file", val);

	//SQL query used to fetch content
	val = wvb_new_string_or_null(wvb_tmpl->content_query);
	json_object_object_add(obj, "content_query", val);

	wvb_query_content(wvb_tmpl);

	if(wvb_tmpl->content == NULL) {
		array = NULL;
		goto ret;
	} else {
		array = json_object_new_array();
	}

	/*
	 * ["string1", "string2", ...],
	 * ["string1", "string2", ...],
	 * ...
	 */

	for(int i = 0; wvb_tmpl->content[i] != NULL; i++) {
		array2 = json_object_new_array();
		for(int x = 0; wvb_tmpl->content[i][x] != NULL; x++) {
			val = wvb_new_string_or_null(wvb_tmpl->content[i][x]);
			json_object_array_add(array2, val);
		}
		json_object_array_add(array, array2);
	}

ret:
	//content array
	json_object_object_add(obj, "content", array);

	return obj;
}

/*
 * Convert a wvb_page into a json object
 */
json_object *wvb_page_to_json_object(WVB_PAGE *wvb_page)
{
	json_object *obj, *array, *val;

	obj = json_object_new_object();

	val = wvb_new_string_or_null(wvb_page->path);
	json_object_object_add(obj, "path", val);

	val = wvb_new_string_or_null(wvb_page->title);
	json_object_object_add(obj, "title", val);

	val = wvb_new_string_or_null(wvb_page->name);
	json_object_object_add(obj, "name", val);

	val = json_object_new_boolean((json_bool) wvb_page->display);
	json_object_object_add(obj, "display", val);

	array = json_object_new_array();

	for(int i = 0; i < wvb_page->template_count; i++) {
		val = wvb_template_to_json_object(&wvb_page->template[i]);
		json_object_array_add(array, val);
	}

	json_object_object_add(obj, "template", array);

	return obj;

}

json_object *wvb_config_to_json_object(WVB_CONFIG *wvb_config)
{
	json_object
		*json,
		*val,
		*page_array;
	int status;

	status = sqlite3_open(wvb_config->database, &db);
	if(db == NULL) {
		LOG_SQL_ERROR(db, status);
		goto err;
	}

	//(re-)create tables and triggers if they were dropped or non-existent
	wvb_database_manage(db);

	json = json_object_new_object();
	val = wvb_new_string_or_null(wvb_config->prefix);

	json_object_object_add(json, "prefix", val);

	json_object_put(val);
	page_array = json_object_new_array();

	for(int i = 0; i < wvb_config->page_count; i++) {
		val = wvb_page_to_json_object(&wvb_config->page[i]);

		if(val == NULL)
			goto err;

		json_object_array_add(page_array, val);
	}

	json_object_object_add(json, "page", page_array);

	goto ret;
err:
	//FIXME properly destruct json objects
	json = NULL;
ret:
	status = sqlite3_close(db);
	LOG_SQL_ERROR(db, status);

	return json;
}
