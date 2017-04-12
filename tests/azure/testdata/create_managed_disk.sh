#!/usr/bin/env nash

import ../../azure/login
import ../../azure/storage

resgroup = $ARGS[1]
location = $ARGS[2]
name     = $ARGS[3]
size     = $ARGS[4]
sku      = $ARGS[5]

azure_login()

disk <= azure_disk_new($name, $resgroup, $location)
disk <= azure_disk_set_size($disk, $size)
disk <= azure_disk_set_sku($disk, $sku)

azure_disk_create($disk)
