#!/usr/bin/env nash

import klb/azure/login
import klb/azure/lb

pool_name = $ARGS[1]
resgroup = $ARGS[2]
lb_name = $ARGS[3]
output = $ARGS[4]

azure_login()

addrpool_id, err <= azure_lb_addresspool_get_id($pool_name, $resgroup, $lb_name)
if $err != "" {
	print("error[%s]\n", $err)
}

echo $addrpool_id > $output
echo "done"
