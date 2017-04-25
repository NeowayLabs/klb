#!/usr/bin/env nash

import klb/azure/login
import klb/azure/group
import klb/azure/nic
import klb/azure/subnet
import klb/azure/vm
import klb/azure/availset
import klb/azure/storage
import klb/azure/disk
import klb/azure/vnet
import klb/azure/nsg
import klb/azure/route
import klb/azure/snapshot
import config.sh

fn addsuffix(name) {
	# Providing true uniqueness with the limits on the names is pretty hard :-)
	s <= head -n1 /dev/urandom | md5sum | tr -dc A-Za-z0-9 | cut -b 1-10

	return $name+"-"+$s
}

fn create_subnet(name, cidr) {
	azure_nsg_create($name, $group, $location)
	azure_subnet_create($name, $group, $vnet, $cidr, $name)
	azure_route_table_create($name, $group, $location)

	hoptype = "Internet"

	route <= azure_route_table_route_new("default", $group, $name, "0.0.0.0/0", $hoptype)

	azure_route_table_route_create($route)
}

fn create_vm(name, subnet) {
	# create ssh key
	accessdir = "/tmp/.config/ssh/"
	accesskey = $accessdir+"id_rsa-"+$name

	-test -e $accesskey

	if $status != "0" {
		mkdir -p $accessdir
		ssh-keygen -f $accesskey -P ""
	}

	# create nic
	nic <= azure_nic_new($name, $group, $location)
	nic <= azure_nic_set_vnet($nic, $vnet)
	nic <= azure_nic_set_subnet($nic, $subnet)

	azure_nic_create($nic)

	vm   <= azure_vm_new($name, $group, $location)
	vm   <= azure_vm_set_vmsize($vm, $vm_size)
	vm   <= azure_vm_set_username($vm, $vm_username)

	nics = ($name)

	vm   <= azure_vm_set_nics($vm, $nics)
	vm   <= azure_vm_set_osdiskname($vm, $name+".vhd")
	vm   <= azure_vm_set_imageurn($vm, $vm_image_urn)
	vm   <= azure_vm_set_publickeyfile($vm, $accesskey+".pub")

	azure_vm_create($vm)
}

azure_login()

if len($ARGS) == "2" {
	group = $ARGS[1]
	
	echo "using existent resource group: "+$group
} else {
	echo "creating new resource group"
	
	azure_group_create($group, $location)
}

echo "creating VNET"

azure_vnet_create($vnet, $group, $location, $vnet_cidr, $vnet_dns_servers)

echo "creating subnet"

create_subnet($subnet_name, $subnet_cidr)

echo "creating virtual machine"

create_vm($vm_name, $subnet_name)
azure_vm_disk_attach_new($vm_name, $group, "premiumDisk", "10", "Premium_LRS")
azure_vm_disk_attach_new($vm_name, $group, "standardDisk", "20", "Standard_LRS")

echo "created main VM"
echo "creating backup VM"

vm_backup_name = $vm_name+"-backup"

create_vm($vm_backup_name, $subnet_name)

echo "getting ID of original VM data disks"

ids <= azure_vm_get_datadisks_ids($vm_name, $group)

echo "generating snapshots from original VM"

for id in $ids {
	snapshot_name <= addsuffix("snapshot")

	echo "creating snapshot: "+$snapshot_name+" from id: "+$id

	snapshotid <= azure_snapshot_create($snapshot_name, $group, $id)

	echo "created snapshot id: "+$snapshotid
	echo "attaching it to backup vm"

	# TODO: attach not working
	# azure_vm_disk_attach($vm_backup_name, $group, $snapshotid)
}

echo
echo "finished with no errors lol"
