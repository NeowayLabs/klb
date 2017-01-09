#!/usr/bin/env nash

import ../../azure/all

resgroup = $ARGS[0]
azure_group_delete($resgroup)
