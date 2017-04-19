#!/usr/bin/env nash

## Resource Group Settings
group    = "vnet-pub-priv"
location = "eastus"

## Vnet Settings
vnet             = "vnet"
vnet_cidr        = "10.50.0.0/16"
vnet_dns_servers = ("8.8.8.8" "8.8.4.4")

subnet_pub_name  = "public"
subnet_pub_cidr  = "10.50.1.0/24"
subnet_priv_name = "private"
subnet_priv_cidr = "10.50.2.0/24"

## VMs Settings

vm_size         = "Basic_A2"
vm_username     = "core"
vm_image_urn    = "CoreOS:CoreOS:Stable:1298.6.0"
vm_storage_type = "LRS"

nat_name        = "nat"
nat_subnet      = "public"
nat_address     = "10.50.1.100"

bastion_name    = "bastion"
bastion_subnet  = "public"
bastion_address = "10.50.1.200"

app_name        = "app"
app_subnet      = "private"
app_address     = "10.50.2.10"
