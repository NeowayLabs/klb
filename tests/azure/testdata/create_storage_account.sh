#!/usr/bin/env nash

import klb/azure/login
import klb/azure/storage

resgroup = $ARGS[1]
storacc  = $ARGS[2]
location = $ARGS[3]
storage_sku = $ARGS[4]

azure_login()

azure_storage_account_create($storacc, $resgroup, $location, "LRS", "Storage")
