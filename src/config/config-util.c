#include <stdio.h>
#include <string.h>

#include <json.h>
#include <sqlite3.h>

#include "config.h"
#include "config-util.h"

#define WVB_INDEXOF_NOTFOUND -1

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
int wvb_indexof_path(WVB_CONFIG *wvb_config, const char *path)
{
	for(int i = 0; i < wvb_config->page_count; i++)
		if(strcmp(path, wvb_config->page[i].path) == 0)
			return i;

	return WVB_INDEXOF_NOTFOUND;
}

int wvb_query_content(WVB_TEMPLATE *wvb_tmpl)
{
	
}

json_object *wvb_template_to_json_object(WVB_TEMPLATE *wvb_tmpl)
{
	json_object *obj, *array, *val;

	obj = json_object_new_object();

	val = wvb_new_string_or_null(wvb_tmpl->name);
	json_object_object_add(obj, "name", val);

	val = wvb_new_string_or_null(wvb_tmpl->file);
	json_object_object_add(obj, "file", val);

	val = wvb_new_string_or_null(wvb_tmpl->content_query);
	json_object_object_add(obj, "content_query", val);

	return obj;

	wvb_query_content(wvb_tmpl);

	//TODO allocate and fill content array
	for(int i = 0; i < wvb_tmpl->content_length; i++) {

	}
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

	return json;
err:
	return NULL;
}
