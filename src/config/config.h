#pragma once

#include <libconfig.h>

typedef struct {
	const char *name;
	const char *file;
	const char *content_query;

	const char ***content;
	int content_length;

	config_setting_t *setting;
} WVB_TEMPLATE;

typedef struct {
	const char *path;
	const char *title;
	const char *name;

	WVB_TEMPLATE *template;

	int template_count;
	int display;
	config_setting_t *setting;
} WVB_PAGE;

typedef struct {
	const char *prefix;

	WVB_PAGE *page;

	int page_count;
	config_t conf;
	config_setting_t *setting;
} WVB_CONFIG;

int wvb_parse_config(const char *file, WVB_CONFIG *wvb_config);
