import klb/azure/group
import klb/azure/storage

# Creates a new azure blob filesystem
#
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
# This is specially useful on read/write scenarios.
#
# If it succeeds you can use the returned instance to list/send/get files from
# Azure BLOB storage.
fn azure_blob_fs_create(resgroup, location, accountname, sku, tier, containername) {
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
	return azure_blob_fs_new($resgroup, $accountname, $containername), ""
}

# Creates a new blob fs instance, just like azure_blob_fs_create but
# it will not attempt to create any resource. Useful for reading operations.
fn azure_blob_fs_new(resgroup, accountname, containername) {
	return ($resgroup $accountname $containername)
}

# Uploads a single file
fn azure_blob_fs_upload(fs, remotepath, localpath) {
	return azure_storage_blob_upload_by_resgroup(
		azure_blob_fs_container($fs),
		azure_blob_fs_account($fs),
		azure_blob_fs_resgroup($fs),
		$remotepath,
		$localpath
	)
}

# Uploads a local dir to a remote dir
fn azure_blob_fs_upload_dir(fs, remotedir, localdir) {
	# WHY: Make code handling results uniform (no relative path handling)
	localdir <= realpath $localdir
	out, status <= tree -if --noreport $localdir
	if $status != "0" {
		return format("error listing directory[%s], output: %s", $localpath, $out)
	}
	all <= split($out, "\n")
	files = ()
	for a in $all {
		_, status <= test -f $a
		if $status == "0" {
			files <= append($files, $a)
		}
	}

	filesprefix <= format("s:^%s::", $localdir)
	for f in $files {
		remotefilename <= echo $f | sed -e $filesprefix
		remotepath = $remotedir + $remotefilename
		err <= azure_blob_fs_upload($fs, $remotepath, $f)
		if $err != "" {
			return $err
		}
	}

	return ""
}

# List all files on the given dir (as far as Azure has dirs on BLOB storage =P)
fn azure_blob_fs_list(fs, remotedir) {
	# WHY: if you take a look at: az storage blob list --help
	# you will see that there is no way to list all files of a blob,
	# only a maximum of 5000 files (there is no kind of start/cursor/index.
	# Azure is definitely the apex of cloud computing =)
	res, err <= azure_storage_blob_list_by_resgroup(
		azure_blob_fs_container($fs),
		azure_blob_fs_account($fs),
		azure_blob_fs_resgroup($fs),
		"5000")
	return $res, $err
}

fn azure_blob_fs_container(fs) {
	return $fs[2]
}

fn azure_blob_fs_account(fs) {
	return $fs[1]
}

fn azure_blob_fs_resgroup(fs) {
	return $fs[0]
}
