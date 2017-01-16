#!/usr/bin/env nash

import ../../azure/login
import ../../azure/group

resgroup = $ARGS[1]
location = $ARGS[2]

azure_login()
echo "creating resource group: " $resgroup " at: " $location
azure_group_create($resgroup, $location)
