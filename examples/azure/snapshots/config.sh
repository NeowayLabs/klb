#!/usr/bin/env nash

## Resource Group Settings
group    = "klb-examples-snapshots"
location = "eastus2"

## Vnet Settings
vnet             = "vnet"
vnet_cidr        = "10.50.0.0/16"
vnet_dns_servers = ("8.8.8.8" "8.8.4.4")
subnet_name      = "network"
subnet_cidr      = "10.50.1.0/24"

## VMs Settings
vm_name         = "snapshots-test-vm"
vm_size         = "Standard_DS4_v2"
vm_username     = "core"
vm_image_urn    = "CoreOS:CoreOS:Stable:1298.6.0"
snapshots_group = $group+"-snapshots"
