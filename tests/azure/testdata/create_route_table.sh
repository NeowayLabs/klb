#!/usr/bin/env nash

import ../../azure/all

name     = $ARGS[1]
resgroup = $ARGS[2]
location = $ARGS[3]

azure_route_table_create($name, $resgroup, $location)