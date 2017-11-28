#!/usr/bin/env nash

import io

import klb/azure/login
import klb/azure/group
import klb/azure/storage

import config

azure_login()

fn aborterr(err, details) {
        if $err != "" {
                io_println("error[%s] %s", $err, $details)
                exit("1")
        }
}

status <= azure_group_exists($group)
if $status != "0" {
        io_println("no resgroup[%s], creating", $group)
        azure_group_create($group, $location)
}

status <= azure_storage_account_exists($account, $group)
if $status != "0" {
        io_println("no storage account[%s], creating", $account)
        err <=  azure_storage_account_create_blob(
            $account,
            $group,
            $location,
            $sku,
            $tier
        )
        aborterr($err, "creating storage account")
}

status <= azure_storage_container_exists_by_resgroup(
        $container,
        $account,
        $group
)

if $status != "0" {
        io_println("no storage container[%s], creating", $container)
        err <= azure_storage_container_create_by_resgroup(
                $container,
                $account,
                $group,
        )
        aborterr($err, "creating storage container")
}

filename <= mktemp

echo "klb example tests" > $filename

echo "uploading file"
err <= azure_storage_blob_upload_by_resgroup(
        $container,
        $account,
        $group,
        $filename,
        $filename
)

aborterr($err, "uploading file[%s]")
rm -f $filename

echo "uploaded file, listing files"
res, err <= azure_storage_blob_list_by_resgroup(
        $container,
        $account,
        $group,
        "100"
)
aborterr($err, "listing files")

io_println("listed files: %s", $res)
