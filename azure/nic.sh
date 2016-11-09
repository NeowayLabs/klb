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

fn azure_vm_set_vnet(instance, vnet) {
	instance <= append($instance, "--subnet-vnet")
	instance <= append($instance, $vnet)

	return $instance
}

fn azure_vm_set_subnet(instance, subnet) {
	instance <= append($instance, "--subnet-name")
	instance <= append($instance, $subnet)

	return $instance
}

fn azure_vm_set_secgrp(instance, secgrp) {
	instance <= append($instance, "--network-security-group-name")
	instance <= append($instance, $secgrp)

	return $instance
}

fn azure_vm_set_ipfw(instance, ipfw) {
	instance <= append($instance, "--enable-ip-forwarding")
	instance <= append($instance, $ipfw)

	return $instance
}

fn azure_nic_create(instance) {
	azure network nic create $instance
}

fn azure_nic_delete(name, group) {
	(
		azure network nic delete 
			--name $name
			--resource-group $group
	)
}
