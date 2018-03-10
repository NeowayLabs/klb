#!/usr/bin/env nash

import klb/azure/login
import klb/azure/blob/fs

resgroup      = $ARGS[1]
accountname   = $ARGS[2]
containername = $ARGS[3]
remotepath    = $ARGS[4]
output        = $ARGS[5]

azure_login()

fs <= azure_blob_fs_new($resgroup, $accountname, $containername)

files, err <= azure_blob_fs_list($fs, $remotepath)
if $err != "" {
	echo $err
	exit("1")
}

echo "got files: "
echo $files
echo "writing on output"
echo $files > $output
