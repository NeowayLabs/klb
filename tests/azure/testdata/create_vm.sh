#!/usr/bin/env nash

import ../../azure/login
import ../../azure/vm

name     = $ARGS[1]
resgroup = $ARGS[2]
location = $ARGS[3]
ostype   = $ARGS[4]

azure_login()
azure_vm_new($name, $resgroup, $location, $ostype)
