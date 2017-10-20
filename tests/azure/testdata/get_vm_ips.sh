#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm

resgroup = $ARGS[1]
vmname   = $ARGS[2]
output   = $ARGS[3]

azure_login()

res <= azure_vm_get_private_ip_addrs($vmname, $resgroup)
echo $res > $output
echo "done"
