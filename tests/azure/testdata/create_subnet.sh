#!/usr/bin/env nash

import ../../azure/all

subnet = $ARGS[1]
resgroup = $ARGS[2]
vnet = $ARGS[3]
location = $ARGS[4]
nsg = $ARGS[5]

azure_subnet_create($subnet, $resgroup, $vnet, $location, $nsg)