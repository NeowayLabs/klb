#!/usr/bin/env nash

import ../../azure/all

resgroup = $ARGS[1]
azure_group_delete($resgroup)
