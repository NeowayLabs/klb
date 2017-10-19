#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm

resgroup = $ARGS[1]
vmname   = $ARGS[2]
ifaceindex = $ARGS[3]
ipindex  = $ARGS[4]
output   = $ARGS[5]

azure_login()

res <= azure_vm_get_ip_address($vmname, $resgroup, $ifaceindex, $ipindex)
echo $res > $output
echo "done"
