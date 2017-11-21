#!/usr/bin/env nash

import klb/azure/login
import klb/azure/storage

resgroup      = $ARGS[1]
accountname   = $ARGS[2]
containername = $ARGS[3]
remotepath    = $ARGS[4]
localpath     = $ARGS[5]

azure_login()

echo "uploading file"
err <= azure_storage_container_blob_upload_by_resgroup($containername, $accountname, $resgroup, $remotepath, $localpath)

if $err != "" {
	echo "error uploading file: " + $err
	exit("1")
}

echo "success"
