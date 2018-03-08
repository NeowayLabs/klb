import klb/azure/group
import klb/azure/storage

# Creates a new azure blob filesystem
# What do we mean with BLOB file system ?
# Well, BLOB storage has no concept of directories and hierarchy.
# Each container is just a flat namespace.
# Representing hierarchical structures can be done using list operations
# with prefix filters, like listing a dir by using "/dir" as the prefix.
#
# The problem is that uploading and downloading directories will not
# work properly, this idea is not supported on the azure cli (az version 2.0.28).
# Uploading a directory (upload-batch) will not create the "dirs" on the
# container, only the basename of files is used. For example, if you upload
# the directory /tmp/test with a file /tmp/test/1 using upload-batch the result
# on the azure container will be a file named "1".
#
# This was the moment when we started to realize that we had been screwed
# by Azure (again). We never tested the download-batch command supposing that
# it will behave equally bad (maintaining only the basename of files).
#
# AWS S3 is also flat, but the aws tools and API's provides a proper directory
# illusion, you can copy directories pretty much the same way you work with
# cp. This blob fs object will try to do the same for you, but on Azure.
#
# The act of creating this FS will trigger the creation of:
# - The resource group
# - The storage account
# - The container
#
# Making it pretty easy to work, but the function may return an error if it is
# unable to create any of the resources.
#
# If it succeeds you can use the returned instance to list/send/get files from
# Azure BLOB storage.
fn azure_blob_fs_new(resgroup, location, accountname, sku, tier, containername) {
	# WHY: for now creating a group that already exists is ok
	# and the function does not return any value, it will abort =/
	# changing it now may break other people.
	azure_group_create($resgroup, $location)

	status <= azure_storage_account_exists($accountname, $resgroup)
	if $status != "0" {
		err <=  azure_storage_account_create_blob(
		    $accountname,
		    $resgroup,
		    $location,
		    $sku,
		    $tier
		)
		if $err != "" {
			return (), $err
		}
        }

        # WHY: creating a container that already exists do not fail on az
        err <= azure_storage_container_create_by_resgroup(
                $containername,
                $accountname,
                $resgroup,
        )
        if $err != "" {
		return (), $err
        }
	return ($resgroup $accountname $containername), ""
}

# Uploads a single file
fn azure_blob_fs_upload(uploader, remotepath, localpath) {
	return azure_storage_blob_upload_by_resgroup(
		azure_blob_fs_container($uploader),
		azure_blob_fs_account($uploader),
		azure_blob_fs_resgroup($uploader),
		$remotepath,
		$localpath
	)
}

# Uploads a dir
fn azure_blob_fs_upload_dir(uploader, localpath) {
	# TODO
}

fn azure_blob_fs_container(uploader) {
	return $uploader[2]
}

fn azure_blob_fs_account(uploader) {
	return $uploader[1]
}

fn azure_blob_fs_resgroup(uploader) {
	return $uploader[0]
}
