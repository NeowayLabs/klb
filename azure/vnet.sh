# Virtual Network related functions

fn azure_vnet_create(name, group, location, address) {
	(
		azure network vnet create 
			--name $name
			--resource-group $group
			--location $location
			--address-prefixes $address
	)
}

fn azure_vnet_delete(name, group) {
	(
		azure network vnet delete 
			--resource-group $group
			--name $name
	)
}
