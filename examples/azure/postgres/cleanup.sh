#!/usr/bin/env nash

import klb/azure/login
import klb/azure/group

group    = "klb-examples-postgres"

azure_login()
azure_group_delete($group)
