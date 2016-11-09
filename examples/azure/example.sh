#!/usr/bin/env nash

import nashlib/all
import klb/azure/all

GROUP    = "klb-tests"
LOCATION = "eastus"
VNET_NAME = "klb-tests-vnet"
NSG_PUBLIC_NAME = "klb-tests-secgroup-public"
NSG_PRIVATE_NAME = "klb-tests-secgroup-private"
ROUTE_TABLE_PUBLIC_NAME = "klb-tests-route-table-public"
ROUTE_TABLE_PRIVATE_NAME = "klb-tests-route-table-private"
SUBNET_PUBLIC_NAME = "klb-tests-subnet-public"
AVAILSET_NAME = "klb-tests-availset-vm1"
STORAGE_ACCOUNT_NAME = "klb-tests-storage-account"
NIC_NAME = "klb-tests-nic-vm1"
PUBLIC_IP_NAME = "klb-tests-public-ip"
VNET_CIDR = "192.168.0.0/16"
ADDRESS_PUBLIC = "192.168.1.0/24"
ADDRESS_PRIVATE = "192.168.2.0/24"
IPNAT = "192.168.1.100"

fn create_prod() {
	# GROUP
	azure_group_create($GROUP, $LOCATION)
	
	# VNET
	azure_vnet_create($VNET, $GROUP, $LOCATION, $VNETCIDR)

	# NSG
	azure_nsg_create($NSGPUBLIC, $GROUP, $LOCATION)
	azure_nsg_create($NSGPRIVATE, $GROUP, $LOCATION)

	# SUBNET
	azure_subnet_create($SUBNETPUBLIC, $GROUP, $VNET, $ADDRESSPUBLIC, $NSGPUBLIC)
	azure_subnet_create("klb-tests-subnet-private", $GROUP, $VNET, $ADDRESSPRIVATE, $NSGPRIVATE)

	# RULE NSG
	azure_nsg_add_rule("klb-tests-subnet-rule-100", $GROUP, $NSGPUBLIC, "Inbound", "100", "*", "192.168.10.65/32", "22", "Allow")

	# ROUTE TABLE
	azure_route_table_create($ROUTETABLEPUBLIC, $GROUP, $LOCATION)
	azure_route_table_create($ROUTETABLEPRIVATE, $GROUP, $LOCATION)

	# ROUTE
	azure_route_table_add_route("klb-tests-route-private", $GROUP, $ROUTETABLEPRIVATE, "0.0.0.0/0", "VirtualAppliance", $IPNAT)

	# PUBLIC IP
	azure_public_ip_create($PUBLIC_IP_NAME, $GROUP, $LOCATION, "Static")

	# AVAILSET
	azure_availset_create($AVAILSET_NAME, $GROUP, $LOCATION)
	
	# NIC
	nic <= azure_nic_new($NIC_NAME, $GROUP, $LOCATION)
	nic <= azure_nic_set_vnet($nic, $VNET_NAME)
	nic <= azure_nic_set_subnet($nic, $SUBNET_PUBLIC_NAME)
	nic <= azure_nic_set_secgrp($nic, $NSG_PUBLIC_NAME)
	nic <= azure_nic_set_ipfw($nic, "true")
	nic <= azure_nic_set_publicip($nic, $PUBLIC_IP_NAME)
	azure_nic_create($nic)

	# STORAGE ACCOUNT
	azure_storage_account_create($STORAGE_ACCOUNT_NAME, $GROUP, $LOCATION, "LRS", "Storage")

	# VM
	vm <= azure_vm_new("klb-tests-vm-nat", $GROUP, $LOCATION, "Linux")
	vm <= azure_vm_set_vmsize($vm, "Basic_A2")
	vm <= azure_vm_set_username($vm, "core")
	vm <= azure_vm_set_availset($vm, $AVAILSET_NAME)
	vm <= azure_vm_set_vnet($vm, $VNET)
	vm <= azure_vm_set_subnet($vm, $SUBNETPUBLIC)
	vm <= azure_vm_set_nic($vm, $NIC_NAME)
	vm <= azure_vm_set_storageaccount($vm, $STORAGE_ACCOUNT_NAME)
	vm <= azure_vm_set_osdiskvhd($vm, "nat-node-1.vhd")
	vm <= azure_vm_set_imageurn($vm, "CoreOS:CoreOS:Stable:899.17.0")
	# vm <= azure_vm_set_datadiskvhd($vm, datadiskvhd)
	# vm <= azure_vm_set_datadisksize($vm, datadisksize)
	# vm <= azure_vm_set_customdata($vm, customdata)
	# vm <= azure_vm_set_publickeyfile($vm, publickeyfile) 
	azure_vm_create($vm)
}

fn delete_prod() {
	azure_group_delete($GROUP)
}

create_prod()
# delete_prod()
