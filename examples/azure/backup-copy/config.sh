#!/usr/bin/env nash

## Resource Group Settings
group    = "klb-examples-backup-copy"
location = "eastus2"

## Vnet Settings
vnet             = "vnet"
vnet_cidr        = "10.50.0.0/16"
vnet_dns_servers = ("8.8.8.8" "8.8.4.4")
subnet_name      = "network"
subnet_cidr      = "10.50.1.0/24"

## VMs Settings
vm_name         = "vm"
vm_size         = "Standard_DS15_v2"
vm_username     = "core"
vm_image_urn    = "CoreOS:CoreOS:Stable:1298.6.0"
vm_disks_count  = "2"
vm_disks_size   = "50"

backup_prefix   = "klb-excp"
backup_vm_name  = "backup-vm"
backup_location = "westus"
backup_group    = $group + "-" + $backup_location
