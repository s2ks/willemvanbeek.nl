#include <stdio.h>
#include <stdlib.h>

#include <libconfig.h>

#include "config.h"

#define TEMPLATE wvb_config->page[i].template[x]
#define PAGE wvb_config->page[i]
int wvb_parse_config(const char *file, WVB_CONFIG *wvb_config)
{
	int count;

	config_init(&wvb_config->conf);

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
		else
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
		for(int x = 0; i < PAGE.template_count; x++) {
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
