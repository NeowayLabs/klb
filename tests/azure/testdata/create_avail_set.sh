#!/usr/bin/env nash

import klb/azure/login
import klb/azure/availset

resgroup = $ARGS[1]
availset = $ARGS[2]
location = $ARGS[3]

azure_login()
azure_availset_create($availset, $resgroup, $location)
