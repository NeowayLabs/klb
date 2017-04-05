#!/usr/bin/env nash

import klb/azure/nic
import klb/azure/subnet
import klb/azure/vm
import klb/azure/public-ip
import klb/azure/availset
import klb/azure/storage
import klb/azure/disk
import map

# londing configs from config.sh
import config.sh

# create vm
fn build_vms_create(name, subnet, address) {
	# create ssh key
	accesskey = ".config/ssh/id_rsa-"+$name

	-test -e $accesskey

	if $status != "0" {
		mkdir -p .config/ssh
		ssh-keygen -f $accesskey -P ""
	}

	# create storage account
	storage_account <= azure_storage_account_create($name, $group, $location, $vm_storage_type, "Storage")

	# create nic
	nic <= azure_nic_new($name, $group, $location)
	nic <= azure_nic_set_vnet($nic, $vnet)
	nic <= azure_nic_set_subnet($nic, $subnet)
	nic <= azure_nic_set_privateip($nic, $address)

	if $subnet == "public" {
		azure_public_ip_create($name, $group, $location, "Static")

		nic <= azure_nic_set_publicip($nic, $name)
		nic <= azure_nic_set_ipfw($nic, "true")
	}

	azure_nic_create($nic)

	# create vm

	vm <= azure_vm_new($name, $group, $location, "Linux")
	vm <= azure_vm_set_vmsize($vm, $vm_size)
	vm <= azure_vm_set_username($vm, $vm_username)
	vm <= azure_vm_set_vnet($vm, $vnet)
	vm <= azure_vm_set_subnet($vm, $subnet)
	vm <= azure_vm_set_nic($vm, $name)
	vm <= azure_vm_set_storageaccount($vm, $storage_account)
	vm <= azure_vm_set_osdiskvhd($vm, $name+".vhd")
	vm <= azure_vm_set_imageurn($vm, $vm_image_urn)
	vm <= azure_vm_set_publickeyfile($vm, $accesskey+".pub")
	vm <= azure_vm_set_disablebootdiagnostics($vm)

	azure_vm_create($vm)
}

build_vms_create($nat_name, $nat_subnet, $nat_address)
build_vms_create($bastion_name, $bastion_subnet, $bastion_address)
build_vms_create($app_name, $app_subnet, $app_address)
