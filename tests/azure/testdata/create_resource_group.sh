#!/usr/bin/env nash

import ../../azure/all

resgroup = $ARGS[0]
location = $ARGS[1]

echo "creating resource group: " $resgroup " at: " $location
azure_group_create($resgroup, $location)
