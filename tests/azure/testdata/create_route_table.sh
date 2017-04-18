#!/usr/bin/env nash

import klb/azure/all

routetable = $ARGS[1]
resgroup   = $ARGS[2]
location   = $ARGS[3]

azure_login()
azure_route_table_create($routetable, $resgroup, $location)
