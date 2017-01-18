#!/usr/bin/env nash

import ../../azure/all

vnet           = $ARGS[1]
subnet         = $ARGS[2]
resgroup       = $ARGS[3]
routetable     = $ARGS[4]

azure_vnet_set_route_table($vnet, $subnet, $resgroup, $routetable)
