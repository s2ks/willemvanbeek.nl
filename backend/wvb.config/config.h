#pragma once

#include <libconfig.h>

#define WVB_LOOKUP_STRING(setting, name, dest) \
	do { \
		if(config_setting_lookup_string(setting, name, dest) == CONFIG_FALSE) { \
			PRINTWARN("\"%s\" is NULL\n", name); \
			*dest = NULL; \
		} \
	} while(0)

typedef struct {
	const char *name;
	const char *file;
	const char *content_query;

	char ***content;
	int content_length; //FIXME unused

	config_setting_t *setting;
} WVB_TEMPLATE;

typedef struct {
	const char *method;
	char **action;
	int count;

	config_setting_t *setting;
} PAGE_ACTION;

typedef struct {
	const char *path;
	const char *title;
	const char *name;
	const char *type;

	PAGE_ACTION page_action;

	WVB_TEMPLATE *template;

	int template_count;
	int display;
	config_setting_t *setting;
} WVB_PAGE;

typedef struct {
	const char *prefix;
	const char *database;

	WVB_PAGE *page;

	int page_count;
	config_t conf;
	config_setting_t *setting;
} WVB_CONFIG;

int wvb_parse_config(const char *file, WVB_CONFIG *wvb_config);
