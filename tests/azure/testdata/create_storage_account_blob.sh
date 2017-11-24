#!/usr/bin/env nash

import klb/azure/login
import klb/azure/storage

resgroup = $ARGS[1]
name     = $ARGS[2]
location = $ARGS[3]
sku      = $ARGS[4]
tier     = $ARGS[5]

azure_login()

err <= azure_storage_account_create_blob($name, $resgroup, $location, $sku, $tier)
if $err != "" {
	echo "error: " + $err
	exit("1")
}
