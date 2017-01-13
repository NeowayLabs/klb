#!/usr/bin/env nash

import ../../azure/all

resgroup = $ARGS[1]
subnet = $ARGS[2]
location = $ARGS[3]
cidr = $ARGS[4]

azure_subnet_create($subnet, $resgroup, $location, $cidr)
