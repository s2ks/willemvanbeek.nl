#include <stdio.h>

#include <json.h>
#include <sqlite3.h>

#include "config.h"
#include "config-util.h"

static inline json_object *json_new_string_or_null(const char *str)
{
	if(str)
		return json_object_new_string(str);
	else
		return NULL;
}

void print_wvb_config_to_json(WVB_CONFIG *wvb_config)
{
	json_object
		*json,
		*val,
		*page_array,
		*templ_array,
		*obj,
		*obj2;

	json = json_object_new_object();
	val = json_new_string_or_null(wvb_config->prefix);

	json_object_object_add(json, "prefix", val);

	json_object_put(val);
	page_array = json_object_new_array();

	for(int i = 0; i < wvb_config->page_count; i++) {
		WVB_PAGE *page = &wvb_config->page[i];

		obj = json_object_new_object();
		val = json_new_string_or_null(page->path);

		json_object_object_add(obj, "path", val);

		val = json_new_string_or_null(page->title);

		json_object_object_add(obj, "title", val);

		val = json_new_string_or_null(page->name);

		json_object_object_add(obj, "name", val);

		//obj2 = json_object_new_object();

		templ_array = json_object_new_array();

		for(int x = 0; x < wvb_config->page[i].template_count; x++) {
			WVB_TEMPLATE *templ = &wvb_config->page[i].template[x];

			obj2 = json_object_new_object();

			val = json_new_string_or_null(templ->name);
			json_object_object_add(obj2, "name", val);

			val = json_new_string_or_null(templ->file);
			json_object_object_add(obj2, "file", val);

			val = json_new_string_or_null(templ->content_query);
			json_object_object_add(obj2, "content_query", val);

			json_object_array_add(templ_array, obj2);

		}
		json_object_object_add(obj, "template", templ_array);
		json_object_array_add(page_array, obj);
		//json_object_put(obj);
	}

	json_object_object_add(json, "page", page_array);

	printf("%s\n", json_object_to_json_string_ext(json,
				JSON_C_TO_STRING_NOSLASHESCAPE |
				JSON_C_TO_STRING_PRETTY_TAB |
				JSON_C_TO_STRING_PRETTY));
}
