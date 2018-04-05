#!/usr/bin/env nash

import klb/azure/login
import klb/azure/blob/fs

resgroup      = $ARGS[1]
accountname   = $ARGS[2]
containername = $ARGS[3]
remotepath    = $ARGS[4]
output        = $ARGS[5]
timeout = "60"

azure_login()

fs <= azure_blob_fs_new($resgroup, $accountname, $containername, $timeout)

dirs, err <= azure_blob_fs_listdir($fs, $remotepath)
if $err != "" {
	echo $err
	exit("1")
}

if len($dirs) == "0" {
	echo "no dir found, exiting"
	exit("0")
}

print("got dirs[%s]\n",$dirs)
echo "writing on output"
echo $dirs > $output
