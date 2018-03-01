import klb/azure/storage

# Creates a new azure blob uploader
# It will create resgroup/account/container as necessary if they do not exist
fn azure_uploader_new(resgroup, location, accountname, sku, tier, containername) {
	# TODO: add resgroup creation
	# TODO: Test when acc already exists
	# TODO: Test when container already exists
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
fn azure_uploader_upload(uploader, remotepath, localpath) {
	return azure_storage_blob_upload_by_resgroup(
		azure_uploader_container($uploader),
		azure_uploader_account($uploader),
		azure_uploader_resgroup($uploader),
		$remotepath,
		$localpath
	)
}

fn azure_uploader_container(uploader) {
	return $uploader[2]
}

fn azure_uploader_account(uploader) {
	return $uploader[1]
}

fn azure_uploader_resgroup(uploader) {
	return $uploader[0]
}
