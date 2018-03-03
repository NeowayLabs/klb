import klb/azure/group
import klb/azure/storage

# Creates a new azure blob uploader
# It will create resgroup/account/container as necessary if they do not exist
fn azure_blob_uploader_new(resgroup, location, accountname, sku, tier, containername) {
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

# Uploads the blob
fn azure_blob_uploader_upload(uploader, remotepath, localpath) {
	return azure_storage_blob_upload_by_resgroup(
		azure_blob_uploader_container($uploader),
		azure_blob_uploader_account($uploader),
		azure_blob_uploader_resgroup($uploader),
		$remotepath,
		$localpath
	)
}

fn azure_blob_uploader_container(uploader) {
	return $uploader[2]
}

fn azure_blob_uploader_account(uploader) {
	return $uploader[1]
}

fn azure_blob_uploader_resgroup(uploader) {
	return $uploader[0]
}
