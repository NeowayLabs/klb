#!/usr/bin/env nash

import ../../azure/login
import ../../azure/vm

resgroup = $ARGS[1]
vmname   = $ARGS[2]
diskname = $ARGS[3]

azure_login()

azure_vm_disk_attach_by_name($vmname, $resgroup, $diskname)
