#!/usr/bin/env nash

import klb/azure/login
import klb/azure/blob/fs

resgroup      = $ARGS[1]
accountname   = $ARGS[2]
containername = $ARGS[3]
remotepath    = $ARGS[4]
outputpath     = $ARGS[5]
timeout = "60"

azure_login()

fs <= azure_blob_fs_new($resgroup, $accountname, $containername, $timeout)

echo "downloading file"
err <= azure_blob_fs_download($fs, $outputpath, $remotepath)

if $err != "" {
	echo "error downloading file: " + $err
	exit("1")
}

echo "success"
