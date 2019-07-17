#pragma once

#include <libconfig.h>

typedef struct {
	config_t conf;
	config_setting_t *setting;
	const char *prefix;

	struct {
		config_setting_t *setting;
		const char *path;
		struct {
			config_setting_t *setting;
			const char *name;
			const char *file_name;
			const char *content_query;
		} *template

		int template_count;
		int display;
	} *page;

	int page_count;
} WVB_CONFIG;

void wvb_parse_config(const char *file);
