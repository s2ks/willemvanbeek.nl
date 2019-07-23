#include <stdio.h>

#include <json.h>

#include "config.h"
#include "config-util.h"

int main(int argc, char **argv) {
	WVB_CONFIG wvb_config;
	json_object *json;
	const char *errstr;

	if(wvb_parse_config("wvb-config/wvb.conf", &wvb_config) == 0) {
		errstr = config_error_text(&wvb_config.conf);
		PRINTERR("error type %d: %s\n",
				config_error_type(&wvb_config.conf),
				errstr ? errstr : "Unknown error");
		if(errstr)
			PRINTERR("file %s line %d\n",
					config_error_file(&wvb_config.conf),
					config_error_line(&wvb_config.conf));
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
