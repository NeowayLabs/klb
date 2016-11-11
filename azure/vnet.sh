# Virtual Network related functions

fn azure_vnet_create(name, group, location, cidr) {
	azure network vnet create --name $name --resource-group $group --location $location --address-prefixes $cidr
}

fn azure_vnet_delete(name, group) {
	azure network vnet delete --resource-group $group --name $name
}
