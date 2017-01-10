#!/usr/bin/env nash

import ../../azure/all

resgroup = $ARGS[0]
availset = $ARGS[1]
location = $ARGS[2]

azure_availset_create($availset, $resgroup, $location)
