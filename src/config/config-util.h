#pragma once

#include <stdio.h>

#include <json.h>
#include "config.h"

#ifdef DEBUG

#define PRINTERR(...) do {\
	fprintf(stderr, "\nERROR: On line %d in file %s in function %s:\n\t", __LINE__, __FILE__, __func__);\
	fprintf(stderr, __VA_ARGS__);\
} while(0)
#define PRINTWARN(...) do {\
	fprintf(stderr, "\nWARNING: On line %d in file %s in function %s:\n\t", __LINE__, __FILE__, __func__);\
	fprintf(stderr, __VA_ARGS__);\
} while(0)

#else

#define PRINTERR(...) do {\
	fprintf(stderr, "ERROR: "); \
	fprintf(stderr, __VA_ARGS__);\
} while(0)
#define PRINTWARN(...) do {} while(0)

#endif
#ifdef VERBOSE

#define PRINTINFO(...) do {\
	fprintf(stderr, "\nINFO: On line %d in file %s in function %s:\n\t", __LINE__, __FILE__, __func__);\
	fprintf(stderr, __VA_ARGS__);\
} while(0)

#else
#define PRINTINFO(...) do {} while(0)
#endif

json_object *wvb_config_to_json_object(WVB_CONFIG *wvb_config);
