#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm

resgroup = $ARGS[1]
vmname   = $ARGS[2]
diskname = $ARGS[3]
size     = $ARGS[4]
sku      = $ARGS[5]
caching  = $ARGS[6]

azure_login()

echo "creating new disk and attaching it"
echo "name: "+$diskname
echo "size gb: "+$size
echo "sku: "+$sku
echo "caching: "+$caching

azure_vm_disk_attach_new($vmname, $resgroup, $diskname, $size, $sku, $caching)

echo "done"
