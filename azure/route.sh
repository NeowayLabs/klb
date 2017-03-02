# Route Table related functions

fn azure_route_table_create(name, group, location) {
	(
		azure network route-table create --name $name --resource-group $group --location $location
	)
}

fn azure_route_table_delete(name, group) {
	(
		azure network route-table delete --name $name --resource-group $group
	)
}

fn azure_route_table_route_new(name, group, routetable, address, hoptype) {
	instance = (
		"--name"
		$name
		"--resource-group"
		$group
		"--route-table-name"
		$routetable
	    "--address-prefix"
		$address
	    "--next-hop-type"
		$hoptype
	)

	return $instance
}

fn azure_route_table_route_set_hop_address(instance, hopaddress) {
	instance <= append($instance, "--next-hop-ip-address")
	instance <= append($instance, $hopaddress)

	return $instance
}

fn azure_route_table_route_create(instance) {
	azure network route-table route create $instance
}

# azure_route_table_add_route Create route in a Route Table
# deprecated: Use the azure_route_table_route_new function
fn azure_route_table_add_route(name, group, routetable, address, hoptype, hopaddress) {
	(
		azure network route-table route create
						--name $name
						--resource-group $group
						--route-table-name $routetable
						--address-prefix $address
						--next-hop-type $hoptype
						--next-hop-ip-address $hopaddress
	)
}

fn azure_route_table_delete_route(name, group, routetable) {
	(
		azure network route-table route delete
						--name $name
						--resource-group $group
						--route-table-name $routetable
	)
}
