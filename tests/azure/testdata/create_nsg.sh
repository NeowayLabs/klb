#!/usr/bin/env nash

import klb/azure/login
import klb/azure/nsg

name     = $ARGS[1]
resgroup = $ARGS[2]
location = $ARGS[3]

azure_login()
azure_nsg_create($name, $resgroup, $location)
