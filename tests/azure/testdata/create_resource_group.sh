#!/usr/bin/env nash

import ../../azure/all

resgroup = $ARGS[1]
location = $ARGS[2]

echo "creating resource group: " $resgroup " at: " $location
azure_group_create($resgroup, $location)
