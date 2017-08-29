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

fn azure_nic_set_subnet_id(instance, subnetid) {
	instance <= append($instance, "--subnet-id")
	instance <= append($instance, $subnetid)

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

fn azure_nic_set_ip_config_name(instance, name) {
        instance <= append($instance, "--ip-config-name")
	instance <= append($instance, $name)

	return $instance
}

fn azure_nic_set_ipfw(instance, ipfw) {
	instance <= append($instance, "--enable-ip-forwarding")
	instance <= append($instance, $ipfw)

	return $instance
}

fn azure_nic_set_lb_address_pool_ids(instance, addresspoolids) {
	fn join(list, sep) {
		out = ""

		for l in $list {
			out = $out+$l+$sep
		}

		out <= echo $out | sed "s/"+$sep+"$//g"

		return $out
	}

	poolids    <= join($addresspoolids, ",")
	instance <= append($instance, "--lb-address-pool-ids")
	instance <= append($instance, $poolids)

	return $instance
}

fn azure_nic_set_lb_inbound_nat_rule_ids(instance, natruleids) {
	fn join(list, sep) {
		out = ""

		for l in $list {
			out = $out+$l+$sep
		}

		out <= echo $out | sed "s/"+$sep+"$//g"

		return $out
	}

	natids   <= join($natruleids, ",")
	instance <= append($instance, "--lb-inbound-nat-rule-ids")
	instance <= append($instance, $natids)

	return $instance
}

fn azure_nic_add_lb_address_pool(name, ipconfig, group, addrpool_id) {
        return _azure_nic_operate_lb_address_pool($name, $ipconfig, $group, $addrpool_id, "add")
}

fn azure_nic_remove_lb_address_pool(name, ipconfig, group, addrpool_id) {
        return _azure_nic_operate_lb_address_pool($name, $ipconfig, $group, $addrpool_id, "remove")
}

fn azure_nic_create(instance) {
	azure network nic create $instance
}

fn azure_nic_delete(name, group) {
	azure network nic delete --name $name --resource-group $group
}

fn _azure_nic_operate_lb_address_pool(name, ipconfig, group, addrpool_id, op) {
        out, status <= (
                az network nic ip-config address-pool $op
                --address-pool $addrpool_id
                --nic-name $name
                --ip-config-name $ipconfig
                --resource-group $group
        )

        if $status != "0" {
                return format(
                    "error: [%s] on [%s] addrpool[%s] on nic[%s] ipconfig[%s] group[%s]",
                    $out,
                    $op,
                    $addrpool_id,
                    $name,
                    $ipconfig,
                    $group
                )
        }

        return ""
}
