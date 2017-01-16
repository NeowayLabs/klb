#!/usr/bin/env nash

import ../../azure/login
import ../../azure/group

resgroup = $ARGS[1]

azure_login()
azure_group_delete($resgroup)
