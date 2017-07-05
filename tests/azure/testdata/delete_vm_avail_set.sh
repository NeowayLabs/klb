#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm

resgroup = $ARGS[1]
availset = $ARGS[2]

azure_login()
azure_vm_availset_delete($availset, $resgroup)
