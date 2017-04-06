#!/usr/bin/env nash

import ../../azure/login
import ../../azure/vm

resgroup = $ARGS[1]
availset = $ARGS[2]

azure_login()
azure_vm_availset_delete($availset, $resgroup)
