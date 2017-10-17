# Storage related functions

# azure_store_account_create creates a new `storage account`.
# `name` is the storage account name
# `group` is the resource group name
# `location` is the azure region
# `sku` is the SKU name (LRS/ZRS/GRS/RAGRS/PLRS)
# `kind` is the account kind (Storage/BlobStorage)
#
# Following the Azure recommendation, this function add a random 3-digit hash
# before your storage account name and return the `storage account` name created
# Ref: https://docs.microsoft.com/en-us/azure/storage/storage-performance-checklist
fn azure_storage_account_create(name, group, location, sku, kind) {
	prefix <= head -n1 /dev/urandom | md5sum | cut -b 1-3

	storagename = $prefix+$name

	(
		azure storage account create --resource-group $group --location $location --sku-name $sku --kind $kind $storagename
	)

	return $storagename
}

# azure_store_account_delete deletes a exit `storage account`.
# `name` is the storage account name
# `group` is the resource group name
fn azure_storage_account_delete(name, group) {
	azure storage account delete --quiet --resource-group $group $name
}

# azure_storage_account_get_keys gets `storage account` keys list.
# `name` is the storage account name
# `group` is the resource group name
#
# This function return a list with storage accout keys
fn azure_storage_account_get_keys(name, group) {
	keys <= azure storage account keys list --resource-group $group $name | grep "key[0-9]" | awk "{print $3}"
	k    <= split($keys, "\n")

	return $k
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
# `quota` is the storage file container quota (in GB)
# `storage account name` is the storage account name
# `storage account key` is the storage account key
#
# Ref: https://docs.microsoft.com/en-us/azure/storage/storage-how-to-use-files-linux
fn azure_storage_container_create(name, quota, storage, storagekey) {
	(
		azure storage container create --container $name --quota $quota --account-name $storage --account-key $storagekey
	)
}

# azure_store_container_delete deletes a exist `storage container`.
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
