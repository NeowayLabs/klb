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

fn azure_storage_account_delete(name, group) {
	azure storage account delete $name
}
