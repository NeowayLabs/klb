# Network Security Group related functions

fn azure_nsg_create(name, group, location) {
	(
		azure network nsg create --name $name --resource-group $group --location $location
	)
}

fn azure_nsg_delete(name, group) {
	(
		azure network nsg delete --name $name --resource-group $group
	)
}

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

fn azure_nsg_rule_set_description(instance, description) {
	instance <= append($instance, "--description")
	instance <= append($instance, $description)

	return $instance
}

fn azure_nsg_rule_set_protocol(instance, protocol) {
	instance <= append($instance, "--protocol")
	instance <= append($instance, $protocol)

	return $instance
}

fn azure_nsg_rule_set_source_address(instance, address) {
	instance <= append($instance, "--source-address-prefix")
	instance <= append($instance, $address)

	return $instance
}

fn azure_nsg_rule_set_source_port(instance, port) {
	instance <= append($instance, "--source-port-range")
	instance <= append($instance, $port)

	return $instance
}

fn azure_nsg_rule_set_destination_address(instance, address) {
	instance <= append($instance, "--destination-address-prefix")
	instance <= append($instance, $address)

	return $instance
}

fn azure_nsg_rule_set_destination_port(instance, port) {
	instance <= append($instance, "--destination-port-range")
	instance <= append($instance, $port)

	return $instance
}

fn azure_nsg_rule_set_access(instance, access) {
	instance <= append($instance, "--access")
	instance <= append($instance, $access)

	return $instance
}

fn azure_nsg_rule_set_direction(instance, direction) {
	instance <= append($instance, "--direction")
	instance <= append($instance, $direction)

	return $instance
}

fn azure_nsg_rule_create(instance) {
	azure network nsg rule create $instance
}

# azure_nsg_rule_update updates an NSG rule.
# `instance` is the name of the instance.
fn azure_nsg_rule_update(instance) {
	az network nsg rule update $instance
}

fn azure_nsg_delete_rule(name, group, nsgname) {
	(
		azure network nsg rule delete
					--name $name
					--resource-group $group
					--nsg-name $nsgname
	)
}

# azure_nsg_get_id will return the NSG ID.
# `name` is the nsg name
# `group` is name of resource group.
fn azure_nsg_get_id(name, group) {
	out, err <= az network nsg show --resource-group $group --name $name --ouput json >[2=]

	if $err != "0" {
		return ""
	}

	nsgid <= echo -n $out | jq -r ".id"

	return $nsgid
}

# azure_nsg_rule_get_id will return the NSG rule ID.
# `nsgname` is the name of the nsg (network security group).
# `name` is the name of the nsg rule
# `group` is name of resource group.
fn azure_nsg_rule_get_id(nsgname, rulename, group) {
	out, err <= (
		az network nsg rule show
					--nsg-name $nsgname
					--name $rulename
					--resource-group $group
					--output json
					>[2=]
	)

	if $err != "0" {
		return ""
	}

	nsgruleid <= echo -n $out | jq -r ".id"

	return $nsgruleid
}
