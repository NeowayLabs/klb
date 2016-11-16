# Virtual Network related functions

fn azure_vnet_create(name, group, location, cidr) {
	azure network vnet create --name $name --resource-group $group --location $location --address-prefixes $cidr
}

fn azure_vnet_delete(name, group) {
	azure network vnet delete --resource-group $group --name $name
}

fn azure_vnet_set_route_table(name, group, vnet, subnet, routetable) {
	azure network vnet subnet set --name $name --resource-group $group --vnet-name $vnet --route-table-name $routetable
}
