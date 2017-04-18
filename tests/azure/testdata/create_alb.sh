#!/usr/bin/env nash

import klb/azure/login
import klb/azure/group
import klb/azure/nsg
import klb/azure/vnet
import klb/azure/subnet
import klb/azure/lb

resgroup              = $ARGS[1]
location              = $ARGS[2]
cidr                  = $ARGS[3]
subnetaddr            = $ARGS[4]
lbname                = $ARGS[5]
frontendip_name       = $ARGS[6]
frontendip_private_ip = $ARGS[7]
addrpoolname          = $ARGS[8]
vnet                  = "alb-tests-vnet"
securitygroup         = "alb-tests-securitygroup"
subnet                = "alb-tests-subnet"
dnsservers            = ("8.8.8.8" "8.8.4.4")

azure_login()
azure_group_create($resgroup, $location)
azure_vnet_create($vnet, $resgroup, $location, $cidr, $dnsservers)
azure_nsg_create($securitygroup, $resgroup, $location)
azure_subnet_create($subnet, $resgroup, $vnet, $subnetaddr, $securitygroup)

# Our main production use case is using subnet id
subnetid <= azure_subnet_get_id($subnet, $resgroup, $vnet)

azure_lb_create($lbname, $resgroup, $location)

frontip <= azure_lb_frontend_ip_new($frontendip_name, $resgroup)
frontip <= azure_lb_frontend_ip_set_lbname($frontip, $lbname)
frontip <= azure_lb_frontend_ip_set_subnet_id($frontip, $subnetid)
frontip <= azure_lb_frontend_ip_set_private_ip($frontip, $frontendip_private_ip)

azure_lb_frontend_ip_create($frontip)
azure_lb_addresspool_create($addrpoolname, $resgroup, $lbname)
