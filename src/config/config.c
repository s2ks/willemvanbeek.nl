#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fnmatch.h>
#include <limits.h>
#include <dirent.h>

#include <libconfig.h>

#include "config.h"

const char **wvb_include(config_t *cfg, const char *include_dir, const char *path, const char **error)
{
	struct dirent *entry;

	char *open_path = NULL, *full_path = NULL, **path_array = NULL;
	char *file_name;
	int path_len = 0;
	int path_array_count = 0;

	DIR *dir = NULL;

	open_path = calloc(1, PATH_MAX + sizeof(char));
	full_path = malloc(PATH_MAX + sizeof(char));

	file_name = strrchr(path, '/');

	if(file_name == NULL)
		file_name = path;
	else
		file_name += sizeof(char); //skip '/' character

	if(*path != '/') {
		if(include_dir) {
			path_len = strlen(include_dir);

			if(path_len > PATH_MAX)
				goto err;

			strcpy(open_path, include_dir);
		}
	}

	path_len = strrchr(path, '/') - path;

	if(path_len > PATH_MAX)
		goto err;

	if(path_len > 0)
		strncat(open_path, path, path_len);

	dir = opendir(open_path);
	if(dir == NULL)
		goto err;

	path_array = malloc(32 * sizeof(*path_array));

	if(path_array == NULL)
		goto err;

	while((entry = readdir(dir)) != NULL) {
		snprintf(full_path, PATH_MAX, "%s/%s", open_path, entry->d_name);

		if(fnmatch(path, full_path, FNM_PATHNAME) != 0)
			continue;

		if(++path_array_count % 32 == 0)
			path_array = realloc(path_array, path_array_count * sizeof(*path_array));

		if(path_array == NULL)
			goto err;

		path_array[path_array_count - 1] = strdup(full_path);
	}

	if(path_array == NULL)
		goto err;

	path_array[path_array_count] = NULL;

err:
	if(dir)
		closedir(dir);

	free(open_path);
	free(full_path);

	return (const char **) path_array;
}

#define TEMPLATE wvb_config->page[i].template[x]
#define PAGE wvb_config->page[i]
int wvb_parse_config(const char *file, WVB_CONFIG *wvb_config)
{
	int count;

	config_init(&wvb_config->conf);

	config_set_include_func(&wvb_config->conf, wvb_include);

	if(config_read_file(&wvb_config->conf, file) == CONFIG_FALSE)
		goto err;

	if(config_lookup_string(&wvb_config->conf, "serve.prefix", &wvb_config->prefix) == CONFIG_FALSE)
		goto err;

	wvb_config->setting = config_lookup(&wvb_config->conf, "serve.urls");

	if(wvb_config->setting == NULL)
		goto err;

	count = config_setting_length(wvb_config->setting);
	wvb_config->page_count = count;

	if(count > 0)
		wvb_config->page = calloc(1, count * sizeof(*wvb_config->page));
	else
		goto err;

	//fill page struct
	for(int i = 0; i < count; i++) {
		config_setting_t *setting;
		const char *val;

		PAGE.setting = config_setting_get_elem(wvb_config->setting, i);

		if(PAGE.setting == NULL)
			goto err;

		setting = PAGE.setting;

		setting = config_setting_lookup(PAGE.setting, "template");

		if(setting == NULL)
			goto err;

		PAGE.template_count = config_setting_length(setting);

		if(PAGE.template_count > 0)
			PAGE.template = calloc(1, PAGE.template_count * sizeof(*PAGE.template));
		else
			goto err;

		//serve.urls.path
		if(config_setting_lookup_string(PAGE.setting, "path", &val))
			PAGE.path = val;
		else
			goto err;

		if(config_setting_lookup_bool(PAGE.setting, "display", &PAGE.display) == CONFIG_FALSE)
			PAGE.display = 1;

		//fill template struct
		for(int x = 0; x < PAGE.template_count; x++) {
			TEMPLATE.setting = config_setting_get_elem(setting, x);

			if(config_setting_lookup_string(TEMPLATE.setting, "name", &val))
				TEMPLATE.name = val;
			else
				goto err;

			if(config_setting_lookup_string(TEMPLATE.setting, "file", &val))
				TEMPLATE.file = val;
			else
				goto err;

			if(config_setting_lookup_string(TEMPLATE.setting, "content_query", &val))
				TEMPLATE.content_query = val;
			else
				TEMPLATE.content_query = NULL;
		}

	}

	return 1;
err:
	return 0;

}

#undef TEMPLATE
#undef PAGE
