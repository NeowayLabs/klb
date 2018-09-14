#!/usr/bin/env nash

import klb/azure/login
import klb/azure/group
import klb/azure/nsg
import klb/azure/vnet
import klb/azure/subnet
import klb/azure/lb

resgroup              = $ARGS[1]
location              = $ARGS[2]
vnet                  = $ARGS[3]
subnet                = $ARGS[4]
lbname                = $ARGS[5]
frontendip_name       = $ARGS[6]
frontendip_private_ip = $ARGS[7]
addrpoolname          = $ARGS[8]
dnsservers            = ("8.8.8.8" "8.8.4.4")

azure_login()

# Our main production use case is using subnet id
print(
	"getting subnet id for resgroup[%s] vnet[%s] subnet[%s]\n",
	$resgroup,
	$vnet,
	$subnet
)
subnetid, err <= azure_subnet_get_id($subnet, $resgroup, $vnet)

if $err != "" {
	print("error[%s] getting subnetid\n", $err)
	exit("1")
}

albtest <= azure_lb_new($lbname, $resgroup, $location)
albtest <= azure_lb_set_subnetid($albtest, $subnetid)
alb_sku = "Standard" 
albtest <= azure_lb_set_sku($albtest, $alb_sku)
frontend_ip_zone = "1"
albtest <= azure_lb_set_frontend_ip_zone($albtest, $frontend_ip_zone)
albtest <= azure_lb_set_private_ip_address($albtest, $frontendip_private_ip)
albtest <= azure_lb_set_backend_pool_name($albtest, $addrpoolname)
albtest <= azure_lb_set_frontend_ip_name($albtest, $frontendip_name)
azure_lb_create($albtest)