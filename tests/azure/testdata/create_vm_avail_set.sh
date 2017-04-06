#!/usr/bin/env nash

import ../../azure/login
import ../../azure/vm

resgroup = $ARGS[1]
availset = $ARGS[2]
location = $ARGS[3]
updatedomain = $ARGS[4]
faultdomain = $ARGS[5]

azure_login()
availset <= azure_vm_availset_new($availset, $resgroup, $location)
availset <= azure_vm_availset_set_updatedomain($availset, $updatedomain)
availset <= azure_vm_availset_set_faultdomain($availset, $faultdomain)
azure_vm_availset_create($availset)
