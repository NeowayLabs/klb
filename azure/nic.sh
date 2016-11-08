# Network Interface Controller related functions

fn azure_nic_create(name, group, location, vnet, subnet, securitygroup, ipforwarding) {
	(
		azure network nic create 
			--name $name
			--resource-group $group
			--location $location
			--subnet-vnet-name $vnet
			--subnet-name $subnet
			--network-security-group-name $securitygroup
			--enable-ip-forwarding $ipforwarding

	)
}

fn azure_nic_delete(name, group) {
	(
		azure network nic delete 
			--name $name
			--resource-group $group
	)
}
