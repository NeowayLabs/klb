#!/usr/bin/env nash

import ../../azure/login
import ../../azure/vnet

name     = $ARGS[1]
resgroup = $ARGS[2]
location = $ARGS[3]
cidr     = $ARGS[4]

azure_login()
azure_vnet_create($name, $resgroup, $location, $cidr)
