#!/usr/bin/env nash

import klb/azure/all

vnet       = $ARGS[1]
subnet     = $ARGS[2]
resgroup   = $ARGS[3]
routetable = $ARGS[4]

azure_login()
azure_vnet_set_route_table($vnet, $subnet, $resgroup, $routetable)
