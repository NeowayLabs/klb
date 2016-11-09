#!/usr/bin/env nash

import nashlib/all
import klb/azure/all

# RESOURCE Settings
RESOURCE_GROUP_NAME          = "klb-tests-2"
LOCATION                     = "eastus"

# VNET Settings
VNET_NAME                    = $RESOURCE_GROUP_NAME+"-vnet"
VNET_ADDRESS_SPACE           = "192.168.0.0/16"

SUBNET_PUBLIC_NAME           = $RESOURCE_GROUP_NAME+"-subnet-public"
SUBNET_PUBLIC_ADDRESS_RANGE  = "192.168.1.0/24"
SUBNET_PRIVATE_NAME          = $RESOURCE_GROUP_NAME+"-subnet-private"
SUBNET_PRIVATE_ADDRESS_RANGE = "192.168.2.0/24"

NSG_PUBLIC_NAME              = $RESOURCE_GROUP_NAME+"-secgroup-public"
NSG_PRIVATE_NAME             = $RESOURCE_GROUP_NAME+"-secgroup-private"

ROUTE_TABLE_PUBLIC_NAME      = $RESOURCE_GROUP_NAME+"-route-table-public"
ROUTE_TABLE_PRIVATE_NAME     = $RESOURCE_GROUP_NAME+"-route-table-private"

# VMs Settings
VM_NAT_NAME                  = $RESOURCE_GROUP_NAME+"-vm-nat"
VM_NAT_ADDRESS               = "192.168.1.100"
VM_NAT_IMAGE_URN             = "CoreOS:CoreOS:Stable:899.17.0"
VM_NAT_KEY_FILE              = $HOME+"/.ssh/id_rsa.pub"
VM_NAT_AVAILSET_NAME         = $RESOURCE_GROUP_NAME+"-availset-nat"
VM_NAT_STORAGE_ACCOUNT_NAME  = "103storagevmnat"
VM_NAT_OS_DISK_VHD_NAME      = $RESOURCE_GROUP_NAME+"-nat-root.vhd"
VM_NAT_NIC_NAME              = $RESOURCE_GROUP_NAME+"-nic-nat"
VM_NAT_PUBLIC_IP_NAME        = $RESOURCE_GROUP_NAME+"-public-ip-nat"

fn create_resource_group() {
	# GROUP
	azure_group_create($RESOURCE_GROUP_NAME, $LOCATION)
}

fn create_virtual_network() {
	# VNET
	azure_vnet_create($VNET_NAME, $RESOURCE_GROUP_NAME, $LOCATION, $VNET_ADDRESS_SPACE)

	# NSG
	azure_nsg_create($NSG_PUBLIC_NAME, $RESOURCE_GROUP_NAME, $LOCATION)
	azure_nsg_create($NSG_PRIVATE_NAME, $RESOURCE_GROUP_NAME, $LOCATION)

	# SUBNET
	azure_subnet_create($SUBNET_PUBLIC_NAME, $RESOURCE_GROUP_NAME, $VNET_NAME, $SUBNET_PUBLIC_ADDRESS_RANGE, $NSG_PUBLIC_NAME)
	azure_subnet_create($SUBNET_PRIVATE_NAME, $RESOURCE_GROUP_NAME, $VNET_NAME, $SUBNET_PRIVATE_ADDRESS_RANGE, $NSG_PRIVATE_NAME)

	# RULE NSG
	azure_nsg_add_inbound_rule("Inbound-ssh-rule-100", $RESOURCE_GROUP_NAME, $NSG_PUBLIC_NAME, "100", "*", "192.168.0.0/16", "22", "Allow")
	azure_nsg_add_outbound_rule("Outbound-ssh-rule-100", $RESOURCE_GROUP_NAME, $NSG_PUBLIC_NAME, "100", "*", "192.168.0.0/16", "22", "Allow")

	# ROUTE TABLE
	azure_route_table_create($ROUTE_TABLE_PUBLIC_NAME, $RESOURCE_GROUP_NAME, $LOCATION)
	azure_route_table_create($ROUTE_TABLE_PRIVATE_NAME, $RESOURCE_GROUP_NAME, $LOCATION)

	# ROUTE
	azure_route_table_add_route($ROUTE_TABLE_PRIVATE_NAME, $RESOURCE_GROUP_NAME, $ROUTE_TABLE_PRIVATE_NAME, "0.0.0.0/0", "VirtualAppliance", $VM_NAT_ADDRESS)
}

fn create_vm_nat() {
	# PUBLIC IP
	azure_public_ip_create($VM_NAT_PUBLIC_IP_NAME, $RESOURCE_GROUP_NAME, $LOCATION, "Static")

	# AVAILSET
	azure_availset_create($VM_NAT_AVAILSET_NAME, $RESOURCE_GROUP_NAME, $LOCATION)

	# NIC
	nic <= azure_nic_new($VM_NAT_NIC_NAME, $RESOURCE_GROUP_NAME, $LOCATION)
	nic <= azure_nic_set_vnet($nic, $VNET_NAME)
	nic <= azure_nic_set_subnet($nic, $SUBNET_PUBLIC_NAME)
	nic <= azure_nic_set_secgrp($nic, $NSG_PUBLIC_NAME)
	nic <= azure_nic_set_ipfw($nic, "true")
	nic <= azure_nic_set_publicip($nic, $VM_NAT_PUBLIC_IP_NAME)
	nic <= azure_nic_set_privateip($nic, $VM_NAT_ADDRESS)
	azure_nic_create($nic)

	# STORAGE ACCOUNT
	azure_storage_account_create($VM_NAT_STORAGE_ACCOUNT_NAME, $RESOURCE_GROUP_NAME, $LOCATION, "LRS", "Storage")

	# VM
	vm <= azure_vm_new($VM_NAT_NAME, $RESOURCE_GROUP_NAME, $LOCATION, "Linux")
	vm <= azure_vm_set_vmsize($vm, "Basic_A2")
	vm <= azure_vm_set_username($vm, "core")
	vm <= azure_vm_set_availset($vm, $VM_NAT_AVAILSET_NAME)
	vm <= azure_vm_set_vnet($vm, $VNET_NAME)
	vm <= azure_vm_set_subnet($vm, $SUBNET_PUBLIC_NAME)
	vm <= azure_vm_set_nic($vm, $VM_NAT_NIC_NAME)
	vm <= azure_vm_set_storageaccount($vm, $VM_NAT_STORAGE_ACCOUNT_NAME)
	vm <= azure_vm_set_osdiskvhd($vm, $VM_NAT_OS_DISK_VHD_NAME)
	vm <= azure_vm_set_imageurn($vm, $VM_NAT_IMAGE_URN)
	# vm <= azure_vm_set_datadiskvhd($vm, datadiskvhd)
	# vm <= azure_vm_set_datadisksize($vm, datadisksize)
	# vm <= azure_vm_set_customdata($vm, customdata)
	vm <= azure_vm_set_publickeyfile($vm, $VM_NAT_KEY_FILE)
	azure_vm_create($vm)
}

fn delete_resource_group() {
	azure_group_delete($RESOURCE_GROUP_NAME)
}

create_resource_group()
create_virtual_network()
create_vm_nat()
delete_resource_group()
