# Storage related functions

fn azure_storage_account_create(name, group, location, sku, kind) {
	(
		azure storage account create
			--resource-group $group
			--location $location
		    --sku-name $sku
		    --kind $kind
		    $name
	)
}

fn azure_storage_account_delete(name, group) {
	azure storage account delete $name
}
