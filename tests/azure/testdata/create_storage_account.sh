#!/usr/bin/env nash

import klb/azure/login
import klb/azure/storage

resgroup = $ARGS[1]
name     = $ARGS[2]
location = $ARGS[3]
sku      = $ARGS[4]

azure_login()

echo "creating storage account: " + $name
err <= azure_storage_account_create_storage($name, $resgroup, $location, $sku)
if $err != "" {
	echo "error: " + $err
	exit("1")
}
echo "success"
