#pragma once

#include <stdio.h>

#include <json.h>
#include "config.h"

#define PRINTERR(...) do {\
	fprintf(stderr, "ERROR: On line %d in file %s:", __LINE__, __FILE__);\
	fprintf(stderr, __VA_ARGS__);\
} while(0)

#ifdef DEBUG
#define PRINTWARN(...) do {\
	fprintf(stderr, "WARNING: On line %d in file %s:", __LINE__, __FILE__);\
	fprintf(stderr, __VA_ARGS__);\
} while(0)
#define PRINTINFO(...) do {\
	fprintf(stderr, "INFO: On line %d in file %s:", __LINE__, __FILE__);\
	fprintf(stderr, __VA_ARGS__);\
} while(0)
#else
#define PRINTWARN(...) do {} while(0)
#define PRINTINFO PRINTWARN
#endif

json_object *wvb_config_to_json_object(WVB_CONFIG *wvb_config);
