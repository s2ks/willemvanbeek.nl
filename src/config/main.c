#include <stdio.h>

#include <json.h>

#include "config.h"
#include "config-util.h"

int main(int argc, char **argv) {
	WVB_CONFIG wvb_config;
	json_object *json;

	const char *errstr;
	const char *config_path;

	if(argc == 0) {
		PRINTERR("Argument required for path to config file\n");
		return -1;
	} else if(argc > 1) {
		PRINTWARN("Too many arguments, expected 1 but got %d\n", argc);
	}

	config_path = (const char *) argv[0];

	if(wvb_parse_config(config_path, &wvb_config) == 0) {
		errstr = config_error_text(&wvb_config.conf);
		PRINTERR("error type %d: %s on line %d in file %s\n",
				config_error_type(&wvb_config.conf),
				errstr ? errstr : "Unknown error",
				config_error_line(&wvb_config.conf),
				config_error_file(&wvb_config.conf));
		return -1;
	}

	json = wvb_config_to_json_object(&wvb_config);

	if(json)
		fprintf(stdout, "%s\n", json_object_to_json_string_ext(json,
					JSON_C_TO_STRING_NOSLASHESCAPE
#ifdef DEBUG
					| JSON_C_TO_STRING_PRETTY
					| JSON_C_TO_STRING_PRETTY_TAB
#endif
					));
	else
		return -1;

	return 0;
}
