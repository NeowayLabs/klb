# Subnet related functions

fn azure_subnet_create(name, group, vnet, address, securitygroup) {
	(
		azure network vnet subnet create
						--name $name
						--resource-group $group
						--vnet-name $vnet
						--address-prefix $address
						--network-security-group-name $securitygroup
	)
}

fn azure_subnet_get_id(name, group, vnet) {
	out, status <= (
		azure network vnet subnet show $group $vnet $name --json
	)

	if $status != "0" {
		return "", $out
	}

	subnetid <= echo $out | jq -r ".id"

	return $subnetid, ""
}

fn azure_subnet_delete(name, group, vnet) {
	(
		azure network vnet subnet delete
						--resource-group $group
						--vnet-name $vnet
						--name $name
	)
}
