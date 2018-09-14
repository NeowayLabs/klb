#!/usr/bin/env nash

## Resource Group Settings
group       = "klb-examples-snapshots-copy"
location    = "eastus2"
sku         = "Standard_LRS"
disk_size   = "10"
other_group = $group+"-copied"
other_location = "westus"