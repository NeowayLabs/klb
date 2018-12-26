#!/usr/bin/env nash

## Resource Group Settings
group    = "klb-examples-loadbalancer"
location = "eastus2"

## Vnet Settings
vnet             = "vnet"
vnet_cidr        = "10.51.0.0/16"
vnet_dns_servers = ("8.8.8.8" "8.8.4.4")
subnet_name      = "network"
subnet_cidr      = "10.51.1.0/24"

## VMs Settings
vm_name        = "loadBalancedVM"
ipconfig_name  = $vm_name + "-ip-config"
vm_size        = "Standard_A0"
vm_username    = "core"
vm_image_urn   = "CoreOS:CoreOS:Stable:1353.7.0"

## Load Balancer Settings
## Our scenario uses load balancers on different resource groups
lb_group = "klb-examples-loadbalancer2"
lb_name = "AwesomeLoadBalancer"
lb_backend_pool_name = "AwesomeBackendAddressPool"
frontendip_name = "AwesomeFrontendIp"
frontendip_private_ip = "10.51.1.100"
