#!/usr/bin/env nash

import klb/azure/login
import klb/azure/group
import klb/azure/nic
import klb/azure/subnet
import klb/azure/vm
import klb/azure/availset
import klb/azure/storage
import klb/azure/vnet
import klb/azure/nsg
import klb/azure/route

import config.sh


fn log(msg) {
	ts <= date "+%T"
	echo $ts + ":" + $msg
}

fn addsuffix(name) {
	# Providing true uniqueness with the limits on the names is pretty hard :-)
	s <= head -n1 /dev/urandom | md5sum | tr -dc A-Za-z0-9 | cut -b 1-10

	return $name+"-"+$s
}

fn create_subnet(group, name, cidr, location) {
	azure_nsg_create($name, $group, $location)
	azure_subnet_create($name, $group, $vnet, $cidr, $name)
	azure_route_table_create($name, $group, $location)

	hoptype = "Internet"

	route <= azure_route_table_route_new("default", $group, $name, "0.0.0.0/0", $hoptype)

	azure_route_table_route_create($route)
}

fn new_vm_base_vm(name, group, subnet, location) {
	# create nic
	nic <= azure_nic_new($name, $group, $location)
	nic <= azure_nic_set_vnet($nic, $vnet)
	nic <= azure_nic_set_subnet($nic, $subnet)

	azure_nic_create($nic)

	echo "created NIC with success"
	echo "creating new VM instance"

	vm   <= azure_vm_new($name, $group, $location)
	vm   <= azure_vm_set_vmsize($vm, $vm_size)

	nics = ($name)

	vm   <= azure_vm_set_nics($vm, $nics)

	echo "returning new VM instance"
	return $vm
}

fn create_vm(name, subnet) {
	# create ssh key
	vm   <= new_vm_base_vm($name, $group, $subnet, $location)
	vm   <= azure_vm_set_osdiskname($vm, $name)
	vm   <= azure_vm_set_imageurn($vm, $vm_image_urn)
	# create ssh key
	accessdir = "/tmp/.config/ssh/"
	accesskey = $accessdir+"id_rsa-"+$name

	-test -e $accesskey

	if $status != "0" {
		mkdir -p $accessdir
		ssh-keygen -f $accesskey -P ""
	}
	vm   <= azure_vm_set_publickeyfile($vm, $accesskey+".pub")
	vm   <= azure_vm_set_username($vm, $vm_username)

	azure_vm_create($vm)
}

azure_login()

echo "creating new resource group"

azure_group_create($group, $location)
azure_group_create($backup_group, $backup_location)

echo "creating VNET"

azure_vnet_create($vnet, $group, $location, $vnet_cidr, $vnet_dns_servers)
azure_vnet_create($vnet, $backup_group, $backup_location, $vnet_cidr, $vnet_dns_servers)

echo "creating subnet"

create_subnet($group, $subnet_name, $subnet_cidr, $location)
create_subnet($backup_group, $subnet_name, $subnet_cidr, $backup_location)

echo "creating virtual machine"

create_vm($vm_name, $subnet_name)

sequence  <= seq "1" $vm_disks_count
range     <= split($sequence, "\n")

print("creating %q disks with size %q\n", $vm_disks_count, $vm_disks_size)

for i in $range {
	azure_vm_disk_attach_new($vm_name, $group, "disk" + $i, $vm_disks_size, "Premium_LRS", "None")
}

echo "stopping VM"
azure_vm_stop($vm_name, $group)

log("starting backup")

backup, err <= azure_vm_backup_create($vm_name, $group, $backup_prefix, "Standard_LRS")
if $err != "" {
	echo $err
	exit("1")
}

log("created backup: "+$backup)
log("now lets create a copy of the backup on another location")
log(format("source location[%s] target location[%s]", $location, $backup_location))

backup_copy <= format("%s-%s", $backup, $backup_location)
err <= azure_vm_backup_copy($backup, $backup_copy, $backup_location, "Standard_LRS")
if $err != "" {
	echo $err
	exit("1")
}

log(format("created backup copy: [%s]", $backup_copy))

echo
echo "creating backup VM"
backupvm <= new_vm_base_vm($backup_vm_name, $backup_group, $subnet_name, $backup_location)
backupvm <= azure_vm_set_ostype($backupvm, "linux")

echo "restoring backup"
azure_vm_backup_recover($backupvm, "Premium_LRS", "None", $backup_copy)
echo "finished with success"

