#!/usr/bin/env nash

import ../../azure/all

name     = $ARGS[1]
resgroup = $ARGS[2]
location = $ARGS[3]

azure_nsg_create($name, $resgroup, $location)