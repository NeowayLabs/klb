#!/usr/bin/env nash

## Resource Group Settings
group    = "klb-examples-backup"
location = "eastus2"

## Vnet Settings
vnet             = "vnet"
vnet_cidr        = "10.50.0.0/16"
vnet_dns_servers = ("8.8.8.8" "8.8.4.4")
subnet_name      = "network"
subnet_cidr      = "10.50.1.0/24"

## VMs Settings
vm_name       = "backup-test-vm"
vm_size       = "Standard_L4s"
vm_username   = "core"
vm_image_urn  = "CoreOS:CoreOS:Stable:1298.6.0"
backup_prefix = "klb-examples-backups"
