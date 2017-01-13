#!/usr/bin/env nash

import ../../azure/all

name     = $ARGS[0]
resgroup = $ARGS[1]
location = $ARGS[2]
cidr     = $ARGS[3]

azure_vnet_create($name, $resgroup, $location, $cidr)
