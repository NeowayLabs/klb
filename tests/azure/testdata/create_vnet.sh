#!/usr/bin/env nash

import ../../azure/all

name     = $ARGS[1]
resgroup = $ARGS[2]
location = $ARGS[3]
cidr     = $ARGS[4]

azure_vnet_create($name, $resgroup, $location, $cidr)
