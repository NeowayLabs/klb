#!/usr/bin/env nash

import klb/azure/login
import klb/azure/group
import klb/azure/nic
import klb/azure/subnet
import klb/azure/vm
import klb/azure/availset
import klb/azure/vnet
import klb/azure/nsg
import klb/azure/route
import config.sh


accessdir = "/tmp/.config/ssh/"
accesskey = $accessdir+"id_rsa-"+$vm_name
accesskeypub = $accesskey+".pub"

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
	_, status <= test -e $accesskey

	if $status != "0" {
		mkdir -p $accessdir
		ssh-keygen -f $accesskey -P ""
	}

	nic <= azure_nic_new($name, $group, $location)
	nic <= azure_nic_set_vnet($nic, $vnet)
	nic <= azure_nic_set_subnet($nic, $subnet)

	azure_nic_create($nic)

	vm   <= azure_vm_new($name, $group, $location)
	vm   <= azure_vm_set_vmsize($vm, $vm_size)
	vm   <= azure_vm_set_username($vm, $vm_username)

	nics = ($name)

	vm <= azure_vm_set_nics($vm, $nics)
	vm <= azure_vm_set_publickeyfile($vm, $accesskeypub)
	vm <= azure_vm_set_osdiskname($vm, $name)
	vm <= azure_vm_set_imageurn($vm, $vm_image_urn)
        vm <= azure_vm_set_storagesku($vm, "Standard_LRS")

	azure_vm_create($vm)
}

azure_login()

echo "creating new resource group"

azure_group_create($group, $location)

echo "creating VNET"

azure_vnet_create($vnet, $group, $location, $vnet_cidr, $vnet_dns_servers)

echo "creating subnet"

create_subnet($subnet_name, $subnet_cidr)

echo "creating virtual machine"

create_vm($vm_name, $subnet_name)

echo "finished with success"
echo "user: " + $vm_username
echo "private key located at: " + $accesskey
