#!/usr/bin/env nash

import ../../azure/all

resgroup = $ARGS[0]
storacc = $ARGS[1]
location = $ARGS[2]

setenv STORAGE_ACCOUNT_NAME <= azure_storage_account_create($storacc, $resgroup, $location, "LRS", "Storage")
