#!/usr/bin/env nash

import klb/azure/login
import klb/azure/storage

resgroup      = $ARGS[1]
accountname   = $ARGS[2]

azure_login()

print("checking if storage account[%s] resgroup[%s] exists\n", $accountname, $resgroup)
status <= azure_storage_account_exists($accountname, $resgroup)
exit($status)
