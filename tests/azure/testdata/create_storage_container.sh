#!/usr/bin/env nash

import klb/azure/login
import klb/azure/storage

resgroup      = $ARGS[1]
accountname   = $ARGS[2]
containername = $ARGS[3]

azure_login()

print("creating container[%s] on storage account[%s]\n", $containername, $accountname)
err <= azure_storage_container_create_by_resgroup($containername, $accountname, $resgroup)

if $err != "" {
	echo "error creating container: " + $err
	exit("1")
}

echo "success"
