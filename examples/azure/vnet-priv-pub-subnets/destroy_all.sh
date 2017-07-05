#!/usr/bin/env nash

import klb/azure/login
import klb/azure/group

# londing configs from config.sh
import config.sh

azure_login()
azure_group_delete($group)
