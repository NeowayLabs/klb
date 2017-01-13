#!/usr/bin/env nash

import ../../azure/all

resgroup = $ARGS[1]
availset = $ARGS[2]
location = $ARGS[3]

azure_availset_create($availset, $resgroup, $location)
