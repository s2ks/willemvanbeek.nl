#include <stdio.h>
#include "config.h"

int main(void) {
	WVB_CONFIG wvb_config;

	if(wvb_parse_config("../../wvb.backend.conf", &wvb_config) == 0) {
		printf("error in file %s on line %d: %s\n",
				config_error_file(&wvb_config.conf),
				config_error_line(&wvb_config.conf),
				config_error_text(&wvb_config.conf));
		return -1;
	}

	printf("prefix: %s\n", wvb_config.prefix);

	for(int i = 0; i < wvb_config.page_count; i++) {
		printf("page %d path: %s\n", i, wvb_config.page[i].path);

		for(int x = 0; x < wvb_config.page[i].template_count; x++) {
			printf("\ttemplate %d name: %s\n", x, wvb_config.page[i].template[x].name);
			printf("\ttemplate %d file: %s\n", x, wvb_config.page[i].template[x].file);
			printf("\ttemplate %d content_query: %s\n", x, wvb_config.page[i].template[x].content_query);
		}

		printf("\tdisplay? %s\n", wvb_config.page[i].display ? "yes" : "no");
	}

	return 0;
}
