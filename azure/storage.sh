# Storage related functions

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
