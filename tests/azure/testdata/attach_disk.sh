#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm

resgroup = $ARGS[1]
vmname   = $ARGS[2]
srcid    = $ARGS[3]

azure_login()
azure_vm_disk_attach($vmname, $resgroup, $srcid)
