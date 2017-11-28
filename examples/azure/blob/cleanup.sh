#!/usr/bin/env nash

import klb/azure/login
import klb/azure/group

import config

azure_login()
azure_group_delete($group)
