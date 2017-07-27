# Virtual Network related functions

fn azure_vnet_create(name, group, location, cidr, dnsservers) {
	fn join(list, sep) {
		out = ""

		for l in $list {
			out = $out+$l+$sep
		}

		out <= echo $out | sed "s/"+$sep+"$//g"

		return $out
	}

	dns <= join($dnsservers, ",")

	azure network vnet create --name $name --resource-group $group --location $location --address-prefixes $cidr --dns-servers $dns
}

fn azure_vnet_delete(name, group) {
	azure network vnet delete --resource-group $group --name $name
}

fn azure_vnet_set_route_table(vnet, subnet, group, routetable) {
	azure network vnet subnet set --name $subnet --resource-group $group --vnet-name $vnet --route-table-name $routetable
}

# azure_vnet_get_id will return the Virtual Network ID.
# `name` is the Virtual Network name.
# `group` is name of resource group.
fn azure_vnet_get_id(name, group) {
	# redirects stderr into stdout

	out, status <= (
		az network vnet show --resource-group $group --name $name --output json
										>[2=1]
	)

	if $status != "0" {
		return "", $out
	}

	vnetid <= echo -n $out | jq -r ".id"

	return $vnetid, ""
}
