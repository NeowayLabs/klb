#!/usr/bin/env nash

import ../../azure/all

resgroup = $ARGS[1]
storacc = $ARGS[2]
location = $ARGS[3]

setenv STORAGE_ACCOUNT_NAME <= azure_storage_account_create($storacc, $resgroup, $location, "LRS", "Storage")
