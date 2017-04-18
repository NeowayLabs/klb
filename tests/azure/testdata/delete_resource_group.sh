#!/usr/bin/env nash

import klb/azure/login
import klb/azure/group

resgroup = $ARGS[1]

azure_login()
azure_group_delete($resgroup)
