#!/usr/bin/env nash

import klb/azure/login
import klb/azure/group
import klb/azure/vnet
import klb/azure/subnet
import klb/azure/route
import klb/azure/nsg

# londing configs from config.sh
import config.sh

azure_login()

# create resource group
azure_group_create($group, $location)

# create vnet
azure_vnet_create($vnet, $group, $location, $vnet_cidr, $vnet_dns_servers)

fn create_subnet(name, cidr, nexthop) {
	azure_nsg_create($name, $group, $location)
	azure_subnet_create($name, $group, $vnet, $cidr, $name)
	azure_route_table_create($name, $group, $location)

	if $nexthop == "Internet" {
		hoptype = "Internet"

		route <= azure_route_table_route_new("default", $group, $name, "0.0.0.0/0", $hoptype)
	} else {
		hoptype = "VirtualAppliance"

		route <= azure_route_table_route_new("default", $group, $name, "0.0.0.0/0", $hoptype)
		route <= azure_route_table_route_set_hop_address($route, $nexthop)
	}

	azure_route_table_route_create($route)
}

# create public subnet
create_subnet($subnet_pub_name, $subnet_pub_cidr, "Internet")

# create private subnet
create_subnet($subnet_priv_name, $subnet_priv_cidr, $nat_address)
