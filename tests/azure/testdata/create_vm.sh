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
custdata = $ARGS[14]
keyfile  = $ARGS[15]

azure_login()
vm <= azure_vm_new($name, $resgroup, $location, $ostype)
vm <= azure_vm_set_vmsize($name, $vmsize)
vm <= azure_vm_set_username($name, $username)
vm <= azure_vm_set_availset($name, $availset)
vm <= azure_vm_set_vnet($name, $vnet)
vm <= azure_vm_set_subnet($name, $subnet)
vm <= azure_vm_set_nic($name, $nic)
vm <= azure_vm_set_storageaccount($name, $storacc)
vm <= azure_vm_set_osdiskvhd($name, $osdisk)
vm <= azure_vm_set_imageurn($name, $imageurn)
vm <= azure_vm_set_customdata($name, $custdata)
vm <= azure_vm_set_publickeyfile($name, $keyfile)
azure_vm_create($name)
