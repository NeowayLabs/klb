#!/usr/bin/env nash

import klb/azure/login
import klb/azure/group

import config.sh

azure_login()
azure_group_delete($group)
azure_group_delete($lb_group)
