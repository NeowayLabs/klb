#!/usr/bin/env nash

import ../../azure/login
import ../../azure/vm

name     = $ARGS[1]
resgroup = $ARGS[2]
location = $ARGS[3]
ostype   = $ARGS[4]
vmsize   = $ARGS[5]
username = $ARGS[6]
availset = $ARGS[7]
vnet     = $ARGS[8]
subnet   = $ARGS[9]
nic      = $ARGS[10]
storacc  = $ARGS[11]
osdisk   = $ARGS[12]
imageurn = $ARGS[13]
#custdata = $ARGS[14]
#keyfile  = $ARGS[15]

azure_login()
vm <= azure_vm_new($name, $resgroup, $location, $ostype)
vm <= azure_vm_set_vmsize($vm, $vmsize)
#vm <= azure_vm_set_username($vm, $username)
#vm <= azure_vm_set_availset($vm, $availset)
#vm <= azure_vm_set_vnet($vm, $vnet)
#vm <= azure_vm_set_subnet($vm, $subnet)
#vm <= azure_vm_set_nic($vm, $nic)
#vm <= azure_vm_set_storageaccount($vm, $storacc)
vm <= azure_vm_set_osdiskvhd($vm, $osdisk)
vm <= azure_vm_set_imageurn($vm, $imageurn)
#vm <= azure_vm_set_customdata($vm, $custdata)
#vm <= azure_vm_set_publickeyfile($vm, $keyfile)
azure_vm_create($vm)
