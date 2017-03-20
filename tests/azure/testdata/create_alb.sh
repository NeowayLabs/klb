#!/usr/bin/env nash

import ../../azure/login
import ../../azure/lb

resgroup = $ARGS[1]
lbname = $ARGS[2]
location = $ARGS[3]

azure_login()
azure_lb_create($lbname, $resgroup, $location)
