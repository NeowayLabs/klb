#!/usr/bin/env nash

import klb/azure/login
import klb/azure/nsg
import klb/azure/vnet
import klb/azure/subnet

subnet   = $ARGS[1]
resgroup = $ARGS[2]
vnet     = $ARGS[3]
addr     = $ARGS[4]
nsg      = $ARGS[5]

azure_login()
azure_subnet_create($subnet, $resgroup, $vnet, $addr, $nsg)
