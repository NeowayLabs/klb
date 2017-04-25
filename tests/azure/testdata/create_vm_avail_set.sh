#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm

resgroup     = $ARGS[1]
name         = $ARGS[2]
location     = $ARGS[3]
updatedomain = $ARGS[4]
faultdomain  = $ARGS[5]

azure_login()

availset <= azure_vm_availset_new($name, $resgroup, $location)
availset <= azure_vm_availset_set_updatedomain($availset, $updatedomain)
availset <= azure_vm_availset_set_faultdomain($availset, $faultdomain)

azure_vm_availset_create($availset)
