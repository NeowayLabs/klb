#!/usr/bin/env nash

import ../../azure/all

resgroup = $ARGS[1]
name     = $ARGS[2]
location = $ARGS[3]
vnet     = $ARGS[4]
subnet   = $ARGS[5]
addrnic  = $ARGS[6]

azure_login()

nic <= azure_nic_new($name, $resgroup, $location)
nic <= azure_nic_set_vnet($nic, $vnet)
nic <= azure_nic_set_subnet($nic, $subnet)
nic <= azure_nic_set_privateip($nic, $addrnic)

azure_nic_create($nic)
