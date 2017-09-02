#!/usr/bin/env nash

import klb/azure/login
import klb/azure/nic

nic_name = $ARGS[1]
ipconfig_name = $ARGS[2]
resgroup = $ARGS[3]
addrpool_id = $ARGS[4]

azure_login()

print("adding lb address pool[%s] to nic[%s]\n", $addrpool_id, $nic_name)
err <= azure_nic_add_lb_address_pool(
    $nic_name,
    $ipconfig_name,
    $resgroup,
    $addrpool_id
)

if $err != "" {
	print("error[%s] adding load balancer address pool to NIC\n", $err)
	exit("1")
}
