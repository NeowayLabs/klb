#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm
import klb/azure/nic

name        = $ARGS[1]
resgroup    = $ARGS[2]
location    = $ARGS[3]
vmsize      = $ARGS[4]
vnet        = $ARGS[5]
subnet      = $ARGS[6]
ostype      = $ARGS[7]
storagesku  = $ARGS[8]
caching     = $ARGS[9]
bkpresgroup = $ARGS[10]

azure_login()

err <= azure_vm_backup_exists($bkpresgroup)

if $err != "" {
	echo "error: "+$err
	
	exit("1")
}

nicname = $name+"-nic"

# create nic
# FIXME: create NIC outside this script, this impairs the capacity
# of trying again this script if the backup recovery fails
nic <= azure_nic_new($nicname, $resgroup, $location)
nic <= azure_nic_set_vnet($nic, $vnet)
nic <= azure_nic_set_subnet($nic, $subnet)

azure_nic_create($nic)

vm   <= azure_vm_new($name, $resgroup, $location)
vm   <= azure_vm_set_vmsize($vm, $vmsize)

nics = ($nicname)

vm   <= azure_vm_set_nics($vm, $nics)
vm   <= azure_vm_set_ostype($vm, $ostype)

echo "creating vm from backup: "+$bkpresgroup

res <= azure_vm_backup_recover($vm, $storagesku, $caching, $bkpresgroup)

if $res != "" {
	echo "error: "+$res
	
	exit("1")
}

echo "done"
