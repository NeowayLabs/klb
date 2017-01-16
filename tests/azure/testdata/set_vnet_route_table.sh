#!/usr/bin/env nash

import ../../azure/all

name           = $ARGS[1]
resgroup       = $ARGS[2]
vnet           = $ARGS[3]
subnet         = $ARGS[4]
routetable     = $ARGS[5]

azure_vnet_set_route_table($name, $resgroup, $vnet, $subnet, $routetable)