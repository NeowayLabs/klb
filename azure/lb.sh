# Azure Load Balancer related functions

# LB functions

fn azure_lb_create(name, group, location) {
	(
		azure network lb create --name $name --resource-group $group --location $location
	)
}

fn azure_lb_delete(name, group) {
	(
		azure network lb delete --name $name --resource-group $group
	)
}

fn azure_lb_get_id(name, group) {
	resp <= (
		azure network lb show --name $name --resource-group $group --json
	)

	id   <= echo $resp | jq -r ".id"

	return $id
}

# FRONTEND IP functions

fn azure_lb_frontend_ip_new(name, group) {
	instance = ("--name" $name "--resource-group" $group)

	return $instance
}

fn azure_lb_frontend_ip_set_lbname(instance, lbname) {
	instance <= append($instance, "--lb-name")
	instance <= append($instance, $lbname)

	return $instance
}

fn azure_lb_frontend_ip_set_private_ip(instance, private_ip) {
	instance <= append($instance, "--private-ip-address")
	instance <= append($instance, $private_ip)

	return $instance
}

fn azure_lb_frontend_ip_set_id_public_ip(instance, idpublicip) {
	instance <= append($instance, "--public-ip-id")
	instance <= append($instance, $idpublicip)

	return $instance
}

fn azure_lb_frontend_ip_set_public_ip_name(instance, publicipname) {
	instance <= append($instance, "--public-ip-name")
	instance <= append($instance, $publicipname)

	return $instance
}

fn azure_lb_frontend_ip_set_subnet_id(instance, subnetid) {
	instance <= append($instance, "--subnet-id")
	instance <= append($instance, $subnetid)

	return $instance
}

fn azure_lb_frontend_ip_set_subnet_name(instance, subnetname) {
	instance <= append($instance, "--subnet-name")
	instance <= append($instance, $subnetname)

	return $instance
}

fn azure_lb_frontend_ip_set_subnet_vnet_name(instance, subnetvnetname) {
	instance <= append($instance, "--subnet-vnet-name")
	instance <= append($instance, $subnetvnetname)

	return $instance
}

fn azure_lb_frontend_ip_create(instance) {
	(azure network lb frontend-ip create $instance)
}

fn azure_lb_frontend_ip_delete(name, group, lbname) {
	(
		azure network lb frontend-ip delete
						--name $name
						--resource-group $group
						--lb-name $lbname
	)
}

# ADDRESS POLL functions

fn azure_lb_addresspool_create(name, group, lbname) {
	out           <= (
		azure network lb address-pool create
						--resource-group $group
						--lb-name $lbname
						--name $name
						--json
						
	)

	addresspoolid <= echo -n $out | jq -r ".id"

	return $addresspoolid
}

fn azure_lb_addresspool_delete(name, group, lbname) {
	(
		azure network lb address-pool delete
						--name $name
						--resource-group $group
						--lb-name $lbname
	)
}

fn azure_lb_addresspool_get_id(addrpoolname, resgroup, lbname) {
	out <= (
		-az network lb address-pool show
						--resource-group $resgroup
						--lb-name $lbname
						--name $addrpoolname
						--output json
						>[2=]
	)

	if $out == "" {
		return ""
	}

	addresspoolid <= echo -n $out | jq -r ".id"

	return $addresspoolid
}

# RULE functions

fn azure_lb_rule_new(name, group) {
	instance = ("--name" $name "--resource-group" $group)

	return $instance
}

fn azure_lb_rule_set_lbname(instance, lbname) {
	instance <= append($instance, "--lb-name")
	instance <= append($instance, $lbname)

	return $instance
}

fn azure_lb_rule_set_protocol(instance, protocol) {
	instance <= append($instance, "--protocol")
	instance <= append($instance, $protocol)

	return $instance
}

fn azure_lb_rule_set_frontendport(instance, frontendport) {
	instance <= append($instance, "--frontend-port")
	instance <= append($instance, $frontendport)

	return $instance
}

fn azure_lb_rule_set_backendport(instance, backendport) {
	instance <= append($instance, "--backend-port")
	instance <= append($instance, $backendport)

	return $instance
}

fn azure_lb_rule_set_enablefloatingip(instance, enablefloatingip) {
	instance <= append($instance, "--enable-floating-ip")
	instance <= append($instance, $enablefloatingip)

	return $instance
}

fn azure_lb_rule_set_frontendipname(instance, frontendipname) {
	instance <= append($instance, "--frontend-ip-name")
	instance <= append($instance, $frontendipname)

	return $instance
}

fn azure_lb_rule_set_addresspoolname(instance, addresspoolname) {
	instance <= append($instance, "--backend-address-pool-name")
	instance <= append($instance, $addresspoolname)

	return $instance
}

fn azure_lb_rule_set_probename(instance, probename) {
	instance <= append($instance, "--probe-name")
	instance <= append($instance, $probename)

	return $instance
}

# Set a client session persistence for LB's rule with the possible values:
# "SourceIP"
# "SourceIPProtocol"
fn azure_lb_rule_set_sessionpersistence(instance, sessionpersistence) {
	instance <= append($instance, "--load-distribution")
	instance <= append($instance, $sessionpersistence)

	return $instance
}

fn azure_lb_rule_create(instance) {
	(azure network lb rule create $instance)
}

fn azure_lb_rule_delete(name, group, lbname) {
	(
		azure network lb rule delete
					--name $name
					--resource-group $group
					--lb-name $lbname
	)
}

# TODO: INBOUND-NAT-RULE functions

# PROBE functions

fn azure_lb_probe_new(name, group) {
	instance = ("--name" $name "--resource-group" $group)

	return $instance
}

fn azure_lb_probe_set_lbname(instance, lbname) {
	instance <= append($instance, "--lb-name")
	instance <= append($instance, $lbname)

	return $instance
}

fn azure_lb_probe_set_protocol(instance, protocol) {
	instance <= append($instance, "--protocol")
	instance <= append($instance, $protocol)

	return $instance
}

fn azure_lb_probe_set_port(instance, port) {
	instance <= append($instance, "--port")
	instance <= append($instance, $port)

	return $instance
}

fn azure_lb_probe_set_interval(instance, interval) {
	instance <= append($instance, "--interval")
	instance <= append($instance, $interval)

	return $instance
}

fn azure_lb_probe_set_count(instance, count) {
	instance <= append($instance, "--count")
	instance <= append($instance, $count)

	return $instance
}

fn azure_lb_probe_set_path(instance, path) {
	instance <= append($instance, "--path")
	instance <= append($instance, $path)

	return $instance
}

fn azure_lb_probe_create(instance) {
	(azure network lb probe create $instance)
}

fn azure_lb_probe_delete(name, group, lbname) {
	(
		azure network lb probe delete
					--name $name
					--resource-group $group
					--lb-name $lbname
	)
}
