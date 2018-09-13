# Storage related functions

# azure_store_account_create_storage creates a new `storage account` of kind Storage.
# `name` is the storage account name
# `group` is the resource group name
# `location` is the azure region
# `sku` is the SKU name
fn azure_storage_account_create_storage(name, group, location, sku) {
    return _azure_storage_account_create_storage($name, $group, $location, $sku, "Storage")
}

# azure_store_account_create_storagev2 creates a new `storage account` of kind StorageV2.
# `name` is the storage account name
# `group` is the resource group name
# `location` is the azure region
# `sku` is the SKU name
fn azure_storage_account_create_storagev2(name, group, location, sku) {
	return _azure_storage_account_create_storage($name, $group, $location, $sku, "StorageV2")
}

# azure_storage_account_exists checks if a storage account exists.
# Returns "0" if it already exists (success), "1" otherwise.
fn azure_storage_account_exists(name, group) {
	output, status <= (az storage account list
		--resource-group $group |
		jq -r ".[].name"
	)
	if $status != "0" {
		return "1"
	}

	accounts <= split($output, "\n")
	for account in $accounts {
		if $account == $name {
			return "0"
		}
	}

	return "1"
}

# azure_store_account_create_blob creates a new `storage account` of kind BlobStorage.
# `name` is the storage account name
# `group` is the resource group name
# `location` is the azure region
# `sku` is the SKU name
# `tier` is the access tier (Hot/Cool)
fn azure_storage_account_create_blob(name, group, location, sku, tier) {
	# WHY: we have two functions because:
	# The access tier used for billing StandardBlob accounts.
	# Cannot be set for StandardLRS, StandardGRS, StandardRAGRS, or
	# PremiumLRS account types. It is required for StandardBlob
        # accounts during creation.  Allowed values: Cool, Hot. 

	output, status <= (az storage account create 
		--name $name
		--resource-group $group
		--location $location
		--sku $sku
		--kind "BlobStorage"
		--access-tier $tier
		>[2=1]
	)
	if $status != "0" {
		return format("error[%s]", $output)
	}
	return ""
}

# azure_store_account_delete deletes a exit `storage account`.
# `name` is the storage account name
# `group` is the resource group name
fn azure_storage_account_delete(name, group) {
	azure storage account delete --quiet --resource-group $group $name
}

# azure_storage_account_get_keys gets all keys of the given storage account
# `name` is the storage account name
# `group` is the resource group name
#
# This function return a list of triples and an error string.
# Each triple is on the form: (keyname, keyvalue, permissions).
# On error the second return value will be a non empty string.
fn azure_storage_account_get_keys(name, group) {
	jsonkeys, status <= az storage account keys list -g $group -n $name --output json
	if $status != "0" {
		return (), format(
			"error[%s] listing keys for account[%s] group[%s]",
			$jsonkeys,
			$name,
			$group,
		)
	}

	keys = ()
	names <= _azure_storage_account_parse_json_list($jsonkeys, ".[].keyName")
	values <= _azure_storage_account_parse_json_list($jsonkeys, ".[].value")
	permissions <= _azure_storage_account_parse_json_list($jsonkeys, ".[].permissions")

	i = "0"
	# Lazy way to iterate
	for _ in $names {
		# WHY: could use name, prefer to initialize triple uniformly
		keys <= append($keys, ($names[$i] $values[$i] $permissions[$i]))
		i <= expr $i "+" "1"
	}

	return $keys, ""
}

# azure_store_share_create creates a new `storage share`.
# `name` is the storage file share name
# `quota` is the storage file share quota (in GB)
# `storage account name` is the storage account name
# `storage account key` is the storage account key
#
# Ref: https://docs.microsoft.com/en-us/azure/storage/storage-how-to-use-files-linux
fn azure_storage_share_create(name, quota, storage, storagekey) {
	(
		azure storage share create --share $name --quota $quota --account-name $storage --account-key $storagekey
	)
}

# azure_store_share_delete deletes a exist `storage share`.
# `name` is the storage file share name
# `storage account name` is the storage account name
# `storage account key` is the storage account key
fn azure_storage_share_delete(name, storage, storagekey) {
	(
		azure storage share delete --quiet
						--share $name
						--account-name $storage
						--account-key $storagekey
	)
}

# azure_store_container_create creates a new `storage container`.
# `name` is the storage file container name
# `storage account name` is the storage account name
# `storage account key` is the storage account key
#
# Returns empty error string on success, non empty error message otherwise.
fn azure_storage_container_create(name, accountname, accountkey) {
	output, status <= (
		az storage container create
			--name $name
			--account-name $accountname
			--account-key $accountkey
		>[2=1]
	)

	if $status != "0" {
		return format(
			"error[%s] creating container[%s] on account name[%s]",
			$output,
			$name,
			$accountname,
		)
	}

	return ""
}


# azure_storage_container_create_from_resgroup is similar to
# azure_storage_container_create, the difference is that the account
# key will be obtained automatically (the first one found with full
# permissions), you only have to care with the account name and
# resource group name (useful when dealing with default keys).
#
# `name` is the storage file container name
# `accountname` is the storage account name
# `resgroup` is the resource group name
#
# Returns empty error string on success, non empty error message otherwise.
fn azure_storage_container_create_by_resgroup(name, accountname, resgroup) {

	accountkey, err <= _azure_storage_account_get_key_value($accountname, $resgroup)
	if $err != "" {
		return $err
	}

	return azure_storage_container_create($name, $accountname, $accountkey)
}

# azure_storage_container_exists checks if a container exists.
# Returns "0" if it already exists (success), "1" otherwise.
fn azure_storage_container_exists(name, accountname, accountkey) {
	output, status <= (
		az storage container exists
			--name $name
			--account-name $accountname
			--account-key $accountkey
			| jq -r ".exists"
		>[2=1]
	)
	if $status != "0" {
		return $status
	}
	if $output == "true" {
		return "0"
	}
	return "1"
}

# azure_storage_container_exists_by_resgroup checks if a container exists.
# Returns "0" if it already exists (success), "1" otherwise.
fn azure_storage_container_exists_by_resgroup(containername, accountname, resgroup) {
	accountkey, err <= _azure_storage_account_get_key_value($accountname, $resgroup)
	if $err != "" {
		return "1"
	}
	return azure_storage_container_exists($containername, $accountname, $accountkey)
}

# azure_storage_blob_list lists all blobs in a container.
# Returns a list and an empty error string on success, otherwise
# it will return an empty list and a non empty error string.
fn azure_storage_blob_list(
	containername,
	accountname,
	accountkey,
	numresults
) {

	output, status <= (
		az storage blob list
			--container-name $containername
			--account-name $accountname
			--account-key $accountkey
			--num-results $numresults
		>[2=1]
	)

	if $status != "0" {
		return (), format("error[%s] listing blobs", $output)
	}

	namesraw, status <= echo $output | jq -r ".[].name" >[2=1]
	if $status != "0" {
		return (), format("error[%s] parsing[%s]", $namesraw, $output)
	}

	names <= split($namesraw, "\n")

	return $names, ""
}

# azure_storage_blob_list_by_resgroup does the same azure_storage_blob_list
# but using the first account key available.
fn azure_storage_blob_list_by_resgroup(
	containername,
	accountname,
	resgroup,
	numresults
) {
	accountkey, err <= _azure_storage_account_get_key_value($accountname, $resgroup)
	if $err != "" {
		return (), $err
	}

	res, err <= azure_storage_blob_list(
		$containername,
		$accountname,
		$accountkey,
		$numresults
	)

	return $res, $err
}

# azure_storage_blob_exists checks if a blob exists.
# Returns "0" if it already exists (success), "1" otherwise.
fn azure_storage_blob_exists(containername, accountname, accountkey, blobpath) {

	output, status <= (
		az storage blob exists
			--container-name $containername
			--account-name $accountname
			--account-key $accountkey
			--name $blobpath
			| jq -r ".exists"
		>[2=1]
	)

	if $status != "0" {
		return $status
	}

	if $output == "true" {
		return "0"
	}

	return "1"
}

# azure_storage_blob_exists_by_resgroup checks if a blob exists.
# Returns "0" if it already exists (success), "1" otherwise.
fn azure_storage_blob_exists_by_resgroup(containername, accountname, resgroup, blobpath) {
	accountkey, err <= _azure_storage_account_get_key_value($accountname, $resgroup)
	if $err != "" {
		return "1"
	}
	return azure_storage_blob_exists(
		$containername,
		$accountname,
		$accountkey,
		$blobpath,
	)
}

fn azure_storage_blob_download(
	containername,
	accountname,
	accountkey,
	remotepath,
	localpath
) {
	output, status <= (az storage blob download
		-c $containername
		--account-name $accountname
		--account-key $accountkey
		-n $remotepath
		-f $localpath
		>[2=1]
	)

	if $status != "0" {
		return format(
			"error downloading file[%s] from container[%s] account[%s]: %s",
			$localpath,
			$containername,
			$accountname,
			$output,
		)
	}

	return ""
}

fn azure_storage_blob_download_by_resgroup(
	containername,
	accountname,
	resgroup,
	remotepath,
	localpath
) {
	accountkey, err <= _azure_storage_account_get_key_value($accountname, $resgroup)
	if $err != "" {
		return $err
	}
	return azure_storage_blob_download(
		$containername,
		$accountname,
		$accountkey,
		$remotepath,
		$localpath
	)
}

fn azure_storage_blob_upload(
	containername,
	accountname,
	accountkey,
	remotepath,
	localpath
) {
	output, status <= (az storage blob upload
		--account-name $accountname
		--account-key $accountkey
		-f $localpath
		-c $containername
		-n $remotepath
		>[2=1]
	)

	if $status != "0" {
		return format(
			"error[%s] uploading file[%s] to container[%s] on account name[%s]",
			$output,
			$localpath,
			$containername,
			$accountname,
		)
	}

	return ""
}

fn azure_storage_blob_upload_by_resgroup(
	containername,
	accountname,
	resgroup,
	remotepath,
	localpath
) {
	accountkey, err <= _azure_storage_account_get_key_value($accountname, $resgroup)
	if $err != "" {
		return $err
	}
	return azure_storage_blob_upload(
		$containername,
		$accountname,
		$accountkey,
		$remotepath,
		$localpath
	)
}

# azure_store_container_delete deletes a `storage container`.
# `name` is the storage file container name
# `storage account name` is the storage account name
# `storage account key` is the storage account key
fn azure_storage_container_delete(name, storage, storagekey) {
	(
		azure storage container delete --quiet
						--container $name
						--account-name $storage
						--account-key $storagekey
	)
}

fn _azure_storage_account_parse_json_list(data, query) {
	valsraw <= echo $data | jq -r $query
	vals <= split($valsraw, "\n")
	return $vals
}

fn _azure_storage_account_get_key_value(accountname, resgroup) {

	keys, err <= azure_storage_account_get_keys($accountname, $resgroup)
	if $err != "" {
		return "", $err
	}

	for key in $keys {
		permissions = $key[2]
		if $permissions == "Full" {
			return $key[1], ""
		}
	}

	return "", format(
		"unable to find account key with full permissions for account[%s] resgroup[%s]",
		$accountname,
		$resgroup,
	)
}

fn _azure_storage_account_create_storage(name, group, location, sku, kind) {
	output, status <= (az storage account create
		--name $name
		--resource-group $group
		--location $location
		--sku $sku
		--kind $kind
		>[2=1]
	)
	if $status != "0" {
		return format("error[%s]", $output)
	}
	return ""
}
