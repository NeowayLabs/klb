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

# azure_subnet_delete deletes a subnet.
# `name` is the name of the network security group.
# `group` is name of resource group.
# `vnet` is name of vnet.
#
# Returns an empty string on success or an error string on failure.
fn azure_subnet_delete(name, group, vnet) {
        out, status <= azure network vnet subnet delete --resource-group $group --vnet-name $vnet --name $name
        if $status != "0" {
                return format("error[%s] deleting nsg[%s] resgroup[%s] vnet[%s]", $out, $name, $group, $vnet)
        }
        return ""
}
