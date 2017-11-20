# Storage related functions

# azure_store_account_create_storage creates a new `storage account` of kind Storage.
# `name` is the storage account name
# `group` is the resource group name
# `location` is the azure region
# `sku` is the SKU name
fn azure_storage_account_create_storage(name, group, location, sku) {
	output, status <= (az storage account create
		--name $name
		--resource-group $group
		--location $location
		--sku $sku
		--kind "Storage"
		>[2=1]
	)
	if $status != "0" {
		return format("error[%s]", $output)
	}
	return ""
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

fn _azure_storage_account_parse_json_list(data, query) {
	valsraw <= echo $data | jq -r $query
	vals <= split($valsraw, "\n")
	return $vals
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
fn azure_storage_container_create(name, storage, storagekey) {
	(
		azure storage container create $name
						--account-name $storage
						--account-key $storagekey
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
