# Subnet related functions

# azure_subnet_create creates a subnet.
# `name` is the name of the network security group.
# `group` is name of resource group.
# `vnet` is name of vnet.
# `address` is the CIDR prefix of the subnet (eg. 10.120.0.0/16)
# `securitygroup` is the NSG of the subnet.
#
# Returns an empty string on success or an error string on failure.
fn azure_subnet_create(name, group, vnet, address, securitygroup) {
	out, status <= (
		azure network vnet subnet create
						--name $name
						--resource-group $group
						--vnet-name $vnet
						--address-prefix $address
						--network-security-group-name $securitygroup
	)
        if $status != "0" {
                return format(
			"error[%s] creating subnet[%s] resgroup[%s] vnet[%s] cidr[%s] nsg[%s]",
			$out,
			$name,
			$group,
			$vnet,
			$address,
			$securitygroup
		)
        }
        return ""
}

# azure_subnet_exists checks if a subnet already exists.
# `name` is the name of the subnet.
# `group` is name of resource group.
#
# Returns "0" on success and "1" otherwise (god I miss booleans =().
fn azure_subnet_exists(name, group, vnet) {
        # TODO: Untested
        out, _ <= az network vnet subnet show --name $name --resource-group $group --vnet-name $vnet
        if $out == "" {
                return "1"
        }
        return "0"
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
        out, status <= azure network vnet subnet delete -q --resource-group $group --vnet-name $vnet --name $name
        if $status != "0" {
                return format(
			"error[%s] deleting subnet[%s] resgroup[%s] vnet[%s]",
			$out,
			$name,
			$group,
			$vnet
		)
        }
        print("azure_subnet_delete: success: %s\n", $out)
        return ""
}
