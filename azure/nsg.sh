# Network Security Group related functions

# azure_nsg_create creates a network security group.
# `name` is the name of the network security group.
# `group` is name of resource group.
# `location` is the Azure Region.
fn azure_nsg_create(name, group, location) {
	(
		az network nsg create --name $name --resource-group $group --location $location
	)
}

# azure_nsg_delete deletes a network security group.
# `name` is the name of the network security group.
# `group` is name of resource group.
fn azure_nsg_delete(name, group) {
	(
		az network nsg delete --name $name --resource-group $group
	)
}

# azure_nsg_rule_new creates a new instance of "network security group rule".
# `name` is the name of the rule.
# `group` is name of resource group.
# `location` is the Azure Region.
# `priority` it the rule priority, between 100 (highest priority)
# and 4096 (lowest priority). Must be unique for each rule in the collection.
fn azure_nsg_rule_new(name, group, nsgname, priority) {
	instance = (
		"--name"
		$name
		"--resource-group"
		$group
		"--nsg-name"
		$nsgname
		"--priority"
		$priority
	)

	return $instance
}

# azure_nsg_rule_set_description sets the rule description.
# `instance` is the instance of network security group rule.
# `description` is the rule description.
fn azure_nsg_rule_set_description(instance, description) {
	instance <= append($instance, "--description")
	instance <= append($instance, $description)

	return $instance
}

# azure_nsg_rule_set_protocol sets the network protocol.
# `instance` is the instance of network security group rule.
# `protocol` is the Network protocol. Allowed values: *, Tcp, Udp. Default: *.
fn azure_nsg_rule_set_protocol(instance, protocol) {
	instance <= append($instance, "--protocol")
	instance <= append($instance, $protocol)

	return $instance
}

# azure_nsg_rule_set_source_address sets the source address.
# `instance` is the instance of network security group rule.
# `address` CIDR prefix or IP range. Use '*' to match all IPs. Can also use
# 'VirtualNetwork', 'AzureLoadBalancer', and 'Internet'.  Default: *.
fn azure_nsg_rule_set_source_address(instance, address) {
	instance <= append($instance, "--source-address-prefix")
	instance <= append($instance, $address)

	return $instance
}

# azure_nsg_rule_set_source_port sets the source port or port range.
# `instance` is the instance of network security group rule.
# `port` is the port or port range between 0-65535. Use '*' to match all ports. Default: *.
fn azure_nsg_rule_set_source_port(instance, port) {
	instance <= append($instance, "--source-port-range")
	instance <= append($instance, $port)

	return $instance
}

# azure_nsg_rule_set_destination_address sets the destination address.
# `instance` is the instance of network security group rule.
# `address` CIDR prefix or IP range. Use '*' to match all IPs. Can also use
# 'VirtualNetwork', 'AzureLoadBalancer', and 'Internet'.  Default: *.
fn azure_nsg_rule_set_destination_address(instance, address) {
	instance <= append($instance, "--destination-address-prefix")
	instance <= append($instance, $address)

	return $instance
}

# azure_nsg_rule_set_destination_port sets the destination port or port range.
# `instance` is the instance of network security group rule.
# `port` is the port or port range between 0-65535. Use '*' to match all ports. Default: 80.
fn azure_nsg_rule_set_destination_port(instance, port) {
	instance <= append($instance, "--destination-port-range")
	instance <= append($instance, $port)

	return $instance
}

# azure_nsg_rule_set_access sets the rule action.
# `instance` is the instance of network security group rule.
# `access` is the action. Allowed values: Allow, Deny. Default: Allow.
fn azure_nsg_rule_set_access(instance, access) {
	instance <= append($instance, "--access")
	instance <= append($instance, $access)

	return $instance
}

# azure_nsg_rule_set_direction sets the rule direction.
# `instance` is the instance of network security group rule.
# `direction.` is the rule direction. Allowed values: Inbound, Outbound. Default: Inbound.
fn azure_nsg_rule_set_direction(instance, direction) {
	instance <= append($instance, "--direction")
	instance <= append($instance, $direction)

	return $instance
}

# azure_nsg_rule_create creates an NSG rule.
# `instance` is the instance of network security group rule.
fn azure_nsg_rule_create(instance) {
	az network nsg rule create $instance
}

# azure_nsg_rule_update updates an NSG rule.
# `instance` is the instance of network security group rule.
fn azure_nsg_rule_update(instance) {
	az network nsg rule update $instance
}

# azure_nsg_delete_rule deletes network security group rule.
# `name` is the network security name
# `group` is name of resource group.
# `nsgname` is the name of the network security group.
fn azure_nsg_delete_rule(name, group, nsgname) {
	(
		az network nsg rule delete
					--name $name
					--resource-group $group
					--nsg-name $nsgname
	)
}

# azure_nsg_get_id will return the network security group ID.
# `name` is the network security group name.
# `group` is name of resource group.
fn azure_nsg_get_id(name, group) {
	# redirects stderr into stdout
	out, err <= (
		az network nsg show --resource-group $group --name $name --ouput json
										>[2=1]
	)

	if $err != "0" {
		return "", $out
	}

	nsgid <= echo -n $out | jq -r ".id"

	return $nsgid, ""
}

# azure_nsg_rule_get_id will return network security group rule ID.
# `name` is the name of the network security group rule.
# `group` is name of resource group.
# `nsgname` is the name of the network security group.
fn azure_nsg_rule_get_id(name, group, nsgname) {
	out, err <= (
		az network nsg rule show
					--name $name
					--resource-group $group
					--nsg-name $nsgname
					--output json
					>[2=1]
	)

	if $err != "0" {
		return "", $out
	}

	nsgruleid <= echo -n $out | jq -r ".id"

	return $nsgruleid, ""
}
