#pragma once

#include <json.h>
#include <sqlite3.h>

json_object *sql_to_json(sqlite3_stmt *stmt);
