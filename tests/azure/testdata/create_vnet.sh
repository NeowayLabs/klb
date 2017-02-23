#!/usr/bin/env nash

import ../../azure/login
import ../../azure/vnet

name     = $ARGS[1]
resgroup = $ARGS[2]
location = $ARGS[3]
cidr     = $ARGS[4]
dnsservers = $ARGS[5]

azure_login()

parseddns <= split($dnsservers, ",")
azure_vnet_create($name, $resgroup, $location, $cidr, $parseddns)
