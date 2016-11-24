# Route Table related functions

fn azure_route_table_create(name, group, location) {
	(
		azure network route-table create
			--name $name
			--resource-group $group
			--location $location
	)
}

fn azure_route_table_delete(name, group) {
	(
		azure network route-table delete
			--name $name
			--resource-group $group
	)
}

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
