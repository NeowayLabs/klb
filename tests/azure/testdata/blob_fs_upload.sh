#!/usr/bin/env nash

import klb/azure/login
import klb/azure/blob/fs

resgroup      = $ARGS[1]
location      = $ARGS[2]
accountname   = $ARGS[3]
sku           = $ARGS[4]
tier          = $ARGS[5]
containername = $ARGS[6]
remotepath    = $ARGS[7]
localpath     = $ARGS[8]
timeout = "60"

azure_login()

echo "uploading file"
fs, err <= azure_blob_fs_create(
	$resgroup,
	$location,
	$accountname,
	$sku,
	$tier,
	$containername,
	$timeout
)
if $err != "" {
	echo "error creating fs: " + $err
	exit("1")
}

echo "uploading file"
err <= azure_blob_fs_upload($fs, $remotepath, $localpath)

if $err != "" {
	echo "error uploading: " + $err
	exit("1")
}

echo "success"
