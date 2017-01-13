#!/usr/bin/env nash

import ../../azure/all

name     = $ARGS[0]
resgroup = $ARGS[1]

azure_vnet_delete($name, $resgroup)
