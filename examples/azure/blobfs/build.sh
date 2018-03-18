#!/usr/bin/env nash

import io

import klb/azure/login
import klb/azure/blob/fs

import config

azure_login()

fn aborterr(err, details) {
        if $err != "" {
                io_println("error[%s] %s", $err, $details)
                exit("1")
        }
}

fs, err <= azure_blob_fs_create($group, $location, $account, $sku, $tier, $container)
aborterr($err, "creating blob fs")

upload_dir <= mktemp -d
io_println("creating files to upload")

# Validate Azure 5000 limit
filescount = "6000"
sequence    <= seq "1" $filescount
range       <= split($sequence, "\n")

for i in $range {
	filename <= format("%s/%s", $upload_dir, $i)
	echo "test" > $filename
}

io_println("uploading %s files", $filescount)
err <= azure_blob_fs_upload_dir($fs, "/test", $upload_dir)
rm -rf $upload_dir
aborterr($err, "uploading dir")


files, err <= azure_blob_fs_list($fs, "/test")
aborterr($err, "listing files")

io_println("listing uploaded files")
for f in $files {
	echo $f
}
io_println("listed [%s] files", len($files))
