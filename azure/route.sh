# Route Table related functions

fn azure_route_table_create(name, group, location) {
	(
		azure network route-table create --name $name --resource-group $group --location $location
	)
}

# azure_route_table_delete will delete the given route table
# `name` is the route table name
# `group` is name of resource group
#
# Returns an empty string on success, an error message otherwise
fn azure_route_table_delete(name, group) {
	out, status <= azure network route-table delete -q --name $name --resource-group $group
	if $status != "0" {
		return format("error[%s] deleting route table[%s] from resgroup[%s]", $out, $name, $group)
	}
	return ""
}

# azure_route_table_get_id will return the Route Table ID.
# `name` is the route table name.
# `group` is name of resource group.
fn azure_route_table_get_id(name, group) {
	# redirects stderr into stdout

	out, status <= (
		az network route-table show --resource-group $group --name $name --output json
											>[2=1]
	)

	if $status != "0" {
		return "", $out
	}

	routetableid <= echo -n $out | jq -r ".id"

	return $routetableid, ""
}

# fn azure_route_table_route_new creates a new instance of "route".
# `name` is the route name.
# `group` is name of resource group.
# `routetable` is route table name.
# `address` is the destination CIDR to which the route applies.
# `hoptype` is the type of Azure hop the packet should be sent to. Allowed values:
# Internet, None, VirtualAppliance, VirtualNetworkGateway, VnetLocal.
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

# fn azure_route_table_route_set_hop_address sets the IP address packets should be
# forwarded to when using the VirtualAppliance hop type.
# `instance` is the name of the route instance.
# `hopaddress` is the Virtualappliance IP address.
fn azure_route_table_route_set_hop_address(instance, hopaddress) {
	instance <= append($instance, "--next-hop-ip-address")
	instance <= append($instance, $hopaddress)

	return $instance
}

# fn azure_route_table_route_update creates a route in a route table.
# `instance` is the name of the route instance.
fn azure_route_table_route_create(instance) {
	az network route-table route create $instance
}

# fn azure_route_table_route_update updates a route in a route table.
# `instance` is the name of the route instance.
fn azure_route_table_route_update(instance) {
	az network route-table route update $instance
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
	out, status <= (
		azure network route-table route delete
						-q
						--name $name
						--resource-group $group
						--route-table-name $routetable
	)
	if $status != "0" {
		return format(
			"error[%s] deleting route[%s] from resgroup[%s] route table[%s]",
			$out,
			$name,
			$group,
			$routetable
		)
	}

	return ""
}

# azure_route_table_route_get_id will return the route ID.
# `name` is the route name in a route table.
# `group` is name of resource group.
# `routetable` is name of route table.
fn azure_route_table_route_get_id(name, group, routetable) {
	# redirects stderr into stdout

	out, status <= (
		az network route-table route show
						--resource-group $group
						--route-table-name $routetable
						--name $name
						--output json
						>[2=1]
	)

	if $status != "0" {
		return "", $out
	}

	routeid <= echo -n $out | jq -r ".id"

	return $routeid, ""
}
