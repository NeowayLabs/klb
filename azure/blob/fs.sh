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
# Downloads a remote dir to a local dir.
# It will not create the remote dir on the local dir, only its
# contents (including other dirs, recursively).
#
# For example, downloading the remote dir /klb to the local dir /test
# will copy all contents of /klb to /test, like /test/file1 and /test/dir1/file1
# but it will not create a "klb" directory inside the local dir.
#
# This is the behavior of Plan9 dircp (http://man.cat-v.org/plan_9/1/tar)
# and it seems more intuitive than the cp -r behavior of creating only the
# base dir of the source path at the target path if it does not exists
# and to copy only the contents you need the hack of regex expansion:
# cp -r /srcdir/* /targetdir.
#
# This function returns an error string if it fails and "" if it succeeds.
fn azure_blob_fs_download_dir(fs, localdir, remotedir) {
	resgroup <= azure_blob_fs_resgroup($fs)
	account <= azure_blob_fs_account($fs)
	accountkey, err <= _azure_storage_account_get_key_value($account, $resgroup)
	if $err != "" {
		return (), $err
	}

	container <= azure_blob_fs_container($fs)
	pattern = $remotedir + "*"
	out, status <= (
		az storage blob download-batch
			--destination $localdir
			--source $container
			--account-name $account
			--account-key $accountkey
	)
	if $status != "0" {
		return format("error[%s] downloading dir[%s] to[%s]", $out, $remotedir, $localdir)
	}
	return ""
}


# Uploads a local dir to a remote dir.
# It will not create the local dir on the remote dir, only its
# contents (including other dirs, recursively).
#
# For example, copying the local dir /klb to the remote dir /test
# will copy all contents of /klb to /test, like /test/file1 and /test/dir1/file1
# but it will not create a "klb" directory inside the remote dir.
#
# This is the behavior of Plan9 dircp (http://man.cat-v.org/plan_9/1/tar)
# and it seems more intuitive than the cp -r behavior of creating only the
# base dir of the source path at the target path if it does not exists
# and to copy only the contents you need the hack of regex expansion:
# cp -r /srcdir/* /targetdir.
#
# This function returns an error string if it fails and "" if it succeeds.
fn azure_blob_fs_upload_dir(fs, remotedir, localdir) {
	# WHY: Make code handling results uniform (no relative path handling)
	localdir <= realpath $localdir
	out, status <= find $localdir
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

	prefixregex <= format("s:^%s::", $localdir)
	for f in $files {
		remotefilename <= echo $f | sed -e $prefixregex
		remotepath = $remotedir + $remotefilename
		remotepath <= realpath -m $remotepath
		err <= azure_blob_fs_upload($fs, $remotepath, $f)
		if $err != "" {
			return $err
		}
	}

	return ""
}

# List all files on the given dir (as far as Azure has dirs on BLOB storage =P)
# It will return a list only with the files located inside of the given dir.
#
# Other dirs inside will not be listed (no interesting way to differentiate dirs
# from files on the returned list).
#
# To list dirs use azure_blob_fs_listdir instead.
fn azure_blob_fs_list(fs, remotedir) {
	if $remotedir == "" {
		return (), "azure_blob_fs_list: error: remote dir MUST not be empty"
	}

	res, err <= _azure_blob_fs_list_prefix($fs, $remotedir)
	if $err != "" {
		return (), $err
	}

	files = ()
	for path in $res {
		dir <= dirname $path
		if $dir == $remotedir {
			files <= append($files, $path)
		}
	}
	return $files, ""
}

# List all dirs on the given dir (as far as Azure has dirs on BLOB storage =P)
# It will return a list only with the dirs located inside of the given dir.
#
# Files will not be listed (no interesting way to differentiate dirs
# from files on the returned list).
#
# To list files use azure_blob_fs_list instead.
fn azure_blob_fs_listdir(fs, remotedir) {
	res, err <= _azure_blob_fs_list_prefix($fs, $remotedir)
	if $err != "" {
		return (), $err
	}

	dirs = ()

	fn repeated_dir(otherdir) {
		for dir in $dirs {
			if $dir == $otherdir {
				return "0"
			}
		}
		return "1"
	}

	fn add_dir(gotdir) {
		if $gotdir == $remotedir {
			return "1"
		}
		parentdir <= dirname $gotdir
		if $parentdir != $remotedir {
			return "1"
		}

		if repeated_dir($gotdir) == "0" {
			return "1"
		}

		return "0"
	}

	# OMG this code is O(n^2) =D
	# Don't u miss hashmaps ? =P
	# Lucky for us nothing is slower than Azure itself (exponential moderfocker)
	for path in $res {
		gotdir <= dirname $path
		if add_dir($gotdir) == "0" {
			dirs <= append($dirs, $gotdir)
		}
	}

	return $dirs, ""
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

fn _azure_blob_fs_list_prefix(fs, prefix) {
	resgroup <= azure_blob_fs_resgroup($fs)
	account <= azure_blob_fs_account($fs)
	accountkey, err <= _azure_storage_account_get_key_value($account, $resgroup)
	if $err != "" {
		return (), $err
	}

	container <= azure_blob_fs_container($fs)
	# WHY: if you take a look at: az storage blob list --help
	# you will see that there is no way to list all files of a blob,
	# only a maximum of 5000 files (there is no kind of start/cursor/index.
	# Azure is definitely the apex of cloud computing =)
	numresults = "5000"

	output, status <= (
		az storage blob list
			--container-name $container
			--account-name $account
			--account-key $accountkey
			--num-results $numresults
			--prefix $prefix
		>[2=1]
	)

	if $status != "0" {
		return (), format("error[%s] listing blobs", $output)
	}

	namesraw, status <= echo $output | jq -r ".[].name" >[2=1]
	if $status != "0" {
		return (), format("error[%s] parsing[%s]", $namesraw, $output)
	}

	if $namesraw == "" {
		return (), ""
	}

	files <= split($namesraw, "\n")

	return $files, ""
}
