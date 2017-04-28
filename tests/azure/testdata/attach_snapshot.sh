#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm
import klb/azure/disk

resgroup = $ARGS[1]
location = $ARGS[2]
vmname   = $ARGS[3]
diskname = $ARGS[4]
disksku  = $ARGS[5]
snapshotid = $ARGS[6]

azure_login()

disk <= azure_disk_new($diskname, $resgroup, $location)
disk <= azure_disk_set_source($disk, $snapshotid)
disk <= azure_disk_set_sku($disk, $disksku)
azure_disk_create($disk)

azure_vm_disk_attach($vmname, $resgroup, $diskname)
