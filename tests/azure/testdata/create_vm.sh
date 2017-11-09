#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm

name     = $ARGS[1]
resgroup = $ARGS[2]
location = $ARGS[3]
vmsize   = $ARGS[4]
username = $ARGS[5]
availset = $ARGS[6]
nic      = ($ARGS[7])
osdisk   = $ARGS[8]
imageurn = $ARGS[9]
keyfile  = $ARGS[10]
sku      = $ARGS[11]
caching  = $ARGS[12]
tags     = ($ARGS[13])

azure_login()

vm <= azure_vm_new($name, $resgroup, $location)
vm <= azure_vm_set_vmsize($vm, $vmsize)
vm <= azure_vm_set_username($vm, $username)
vm <= azure_vm_set_availset($vm, $availset)
vm <= azure_vm_set_nics($vm, $nic)
vm <= azure_vm_set_osdiskname($vm, $osdisk)
vm <= azure_vm_set_imageurn($vm, $imageurn)
vm <= azure_vm_set_publickeyfile($vm, $keyfile)
vm <= azure_vm_set_storagesku($vm, $sku)
vm <= azure_vm_set_tags($vm, $tags)
if $caching != "None" {
	# OS disk do not support no caching =/
	vm <= azure_vm_set_osdisk_caching($vm, $caching)
}
vm <= azure_vm_set_datadisk_caching($vm, $caching)

azure_vm_create($vm)
