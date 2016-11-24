# Network Interface Controller related functions

fn azure_nic_new(name, group, location) {
	instance = (
		"--name"
		$name
		"--resource-group"
		$group
		"--location"
		$location
	)

	return $instance
}

fn azure_nic_set_vnet(instance, vnet) {
	instance <= append($instance, "--subnet-vnet-name")
	instance <= append($instance, $vnet)

	return $instance
}

fn azure_nic_set_subnet(instance, subnet) {
	instance <= append($instance, "--subnet-name")
	instance <= append($instance, $subnet)

	return $instance
}

fn azure_nic_set_privateip(instance, privateip) {
	instance <= append($instance, "--private-ip-address")
	instance <= append($instance, $privateip)

	return $instance
}

fn azure_nic_set_publicip(instance, publicip) {
	instance <= append($instance, "--public-ip-name")
	instance <= append($instance, $publicip)

	return $instance
}

fn azure_nic_set_secgrp(instance, secgrp) {
	instance <= append($instance, "--network-security-group-name")
	instance <= append($instance, $secgrp)

	return $instance
}

fn azure_nic_set_ipfw(instance, ipfw) {
	instance <= append($instance, "--enable-ip-forwarding")
	instance <= append($instance, $ipfw)

	return $instance
}

fn azure_nic_create(instance) {
	azure network nic create $instance
}

fn azure_nic_delete(name, group) {
	(
		azure network nic delete --name $name --resource-group $group
	)
}
