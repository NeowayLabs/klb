#!/usr/bin/env nash

import ../../azure/all

routetable     = $ARGS[1]
route          = $ARGS[2]
resgroup       = $ARGS[3]
location       = $ARGS[4]
address        = $ARGS[5]
hoptype        = $ARGS[6]
hopaddress     = $ARGS[7]

azure_route_table_create($routetable, $resgroup, $location)
azure_route_table_add_route($route, $resgroup, $routetable, $hoptype, $hopaddress)