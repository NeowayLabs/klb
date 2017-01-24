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
        resp <= (
                azure network vnet subnet show $group $vnet $name --json
        )
        id <= echo $resp | jq ".id"
        return $id
}

fn azure_subnet_delete(name, group, vnet) {
	(
		azure network vnet subnet delete
						--resource-group $group
						--vnet-name $vnet
						--name $name
	)
}
