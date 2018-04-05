#!/usr/bin/env nash

import klb/azure/login
import klb/azure/blob/fs

resgroup      = $ARGS[1]
accountname   = $ARGS[2]
containername = $ARGS[3]
localdir      = $ARGS[4]
remotedir     = $ARGS[5]
timeout = "60"

azure_login()

echo "downloading dir"
fs <= azure_blob_fs_new($resgroup, $accountname, $containername, $timeout)

echo "downloading directory"
err <= azure_blob_fs_download_dir($fs, $localdir, $remotedir)

if $err != "" {
	echo "error downloading: " + $err
	exit("1")
}

echo "success"
