#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fnmatch.h>
#include <limits.h>
#include <dirent.h>
#include <errno.h>

#include <libconfig.h>

#include "config.h"
#include "config-util.h"
#include "database.h"

static const char *wvb_include_err[] = {
	"Error parsing include",
};

/*
 * TODO directory globbing
 */
static const char **wvb_include(config_t *cfg, const char *include_dir, const char *path, const char **error)
{
	struct dirent *entry;

	char *open_path = NULL, *full_path = NULL, *search_path = NULL, **path_array = NULL;
	char *file_name;
	int path_len = 0;
	int path_array_count = 0;

	DIR *dir = NULL;

	open_path = calloc(1, PATH_MAX + sizeof(char));
	full_path = malloc(PATH_MAX + sizeof(char));
	search_path = malloc(PATH_MAX + sizeof(char));

	if(!open_path || !full_path || !search_path) {
		PRINTERR("allocation failed with error code %d\n", errno);
		goto err;
	}

	file_name = strrchr(path, '/');

	if(file_name == NULL)
		file_name = path;
	else
		file_name += sizeof(char); //skip '/' character

	if(*path != '/') {
		if(include_dir) {
			path_len = strlen(include_dir);

			if(path_len > PATH_MAX) {
				PRINTERR("path_len > PATH_MAX (%d)\n", PATH_MAX);
				goto err;
			}

			strcpy(open_path, include_dir);
			strcpy(search_path, include_dir);
		}
	}

	path_len = strrchr(path, '/') - path;

	if(path_len > PATH_MAX) {
		PRINTERR("path_len > PATH_MAX (%d)\n", PATH_MAX);
		goto err;
	}

	if(path_len > 0)
		strncat(open_path, path, path_len);

	strcat(search_path, path);
	dir = opendir(open_path);
	if(dir == NULL) {
		PRINTERR("failed to open directory with path pointed to by %p\n", open_path);
		goto err;
	}

	path_array = malloc(32 * sizeof(*path_array));

	if(path_array == NULL) {
		PRINTERR("allocation failed with error code %d\n", errno);
		goto err;
	}

	while((entry = readdir(dir)) != NULL) {
		snprintf(full_path, PATH_MAX, "%s/%s", open_path, entry->d_name);

		if(entry->d_type != DT_REG)
			continue;

		if(fnmatch(search_path, full_path, FNM_PATHNAME) != 0)
			continue;

		if(++path_array_count % 32 == 0)
			path_array = realloc(path_array, path_array_count * sizeof(*path_array));

		if(path_array == NULL) {
			PRINTERR("reallocation failed with error code %d\n", errno);
			goto err;
		}

		path_array[path_array_count - 1] = strdup(full_path);
	}

	if(path_array == NULL) {
		PRINTERR("path_array = NULL\n");
		goto err;
	}

	path_array[path_array_count] = NULL;

	goto ret;
err:
	*error = wvb_include_err[0];
ret:
	if(dir)
		closedir(dir);

	free(open_path);
	free(full_path);
	free(search_path);

	return (const char **) path_array;
}


//FIXME defines
//TODO clean up
#define TEMPLATE wvb_config->page[i].template[x]
#define PAGE wvb_config->page[i]
int wvb_parse_config(const char *file, WVB_CONFIG *wvb_config)
{
	static char *include = NULL;
	int count, ret = 1;
	char *c;

	c = strrchr(file, '/');
	c += sizeof(char); //include '/' character

	//FIXME destroy after use
	config_init(&wvb_config->conf);

	if(c != NULL) {
		count = c - file;

		if(include)
			free(include);

		include = malloc(count + sizeof(char));
		snprintf(include, count + sizeof(char), "%s", file);

		//FIXME memory leak
		config_set_include_dir(&wvb_config->conf, include);
	}

	config_set_include_func(&wvb_config->conf, wvb_include);

	//FIXME memory leak
	if(config_read_file(&wvb_config->conf, file) == CONFIG_FALSE) {
		PRINTERR("config_read_file returned CONFIG_FALSE\n");
		goto err;
	}

	if(config_lookup_string(&wvb_config->conf, "serve.prefix", &wvb_config->prefix) == CONFIG_FALSE)
		wvb_config->prefix = NULL;

	if(config_lookup_string(&wvb_config->conf, "serve.database", &wvb_config->database) == CONFIG_FALSE) {
		PRINTERR("Failed to look up \"serve.database\"\n");
		goto err;
	}

	wvb_config->setting = config_lookup(&wvb_config->conf, "serve.urls");

	if(wvb_config->setting == NULL){
		PRINTERR("failed to look up \"serve.urls\"\n");
		goto err;
	}

	count = config_setting_length(wvb_config->setting);
	wvb_config->page_count = count;

	if(count > 0) {
		//FIXME memory leak
		wvb_config->page = calloc(1, count * sizeof(*wvb_config->page));
	} else {
		PRINTERR("page count < 1\n");
		goto err;
	}

	//fill page struct
	for(int i = 0; i < count; i++) {
		config_setting_t *setting;
		const char *val;

		PAGE.setting = config_setting_get_elem(wvb_config->setting, i);

		if(PAGE.setting == NULL) {
			PRINTERR("failed to get page %d\n", i);
			goto err;
		}

		setting = PAGE.setting;

		setting = config_setting_lookup(PAGE.setting, "template");

		if(setting == NULL) {
			PRINTERR("failed to get template array for page %d\n", i);
			goto err;
		}

		PAGE.template_count = config_setting_length(setting);

		if(PAGE.template_count > 0) {
			PAGE.template = calloc(1, PAGE.template_count * sizeof(*PAGE.template));
		} else {
			PRINTERR("page %d template count < 1\n", i);
			goto err;
		}

		//serve.urls.path
		if(config_setting_lookup_string(PAGE.setting, "path", &val)) {
			PAGE.path = val;
		} else {
			PRINTWARN("page %d path = NULL\n", i);
			PAGE.path = NULL;
		}

		WVB_LOOKUP_STRING(PAGE.setting, "title", &PAGE.title);
		WVB_LOOKUP_STRING(PAGE.setting, "name", &PAGE.name);

		if(config_setting_lookup_bool(PAGE.setting, "display", &PAGE.display) == CONFIG_FALSE)
			PAGE.display = 1;

		//fill template struct
		for(int x = 0; x < PAGE.template_count; x++) {
			TEMPLATE.setting = config_setting_get_elem(setting, x);

			if(config_setting_lookup_string(TEMPLATE.setting, "name", &val)) {
				TEMPLATE.name = val;
			} else {
				PRINTINFO("page %d template %d name = NULL\n", i, x);
				TEMPLATE.name = NULL;
			}

			if(config_setting_lookup_string(TEMPLATE.setting, "file", &val)) {
				TEMPLATE.file = val;
			} else {
				PRINTINFO("page %d template %d file = NULL\n", i, x);
				TEMPLATE.file = NULL;
			}

			if(config_setting_lookup_string(TEMPLATE.setting, "content_query", &val)) {
				TEMPLATE.content_query = val;
			} else {
				PRINTINFO("page %d template %d content_query = NULL\n", i, x);
				TEMPLATE.content_query = NULL;
			}
		}

	}

	goto ret;

err:
	ret = 0;
ret:
	//TODO config should be destroyed, but uncommenting this causes problems (check valgrind)
	//config_destroy(&wvb_config->conf);
	return ret;

}

#undef TEMPLATE
#undef PAGE
