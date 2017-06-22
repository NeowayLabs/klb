# Public IP related functions

fn azure_public_ip_create(name, group, location, allocation) {
	(
		azure network public-ip create --name $name --resource-group $group --location $location --allocation-method $allocation
	)
}

fn azure_public_ip_get_address(name, group) {

	public_ip_info    <= azure network public-ip show -g $group -n $name --json
	if $public_ip_info == "{}" {
	   print("no public ip found for resource group: '" + $group + "' and name: '" + $name + "'\n")
	   return ""
	}
	public_ip_address <= echo $public_ip_info | jq -r ".ipAddress"

	return $public_ip_address
}

fn azure_public_ip_delete(name, group) {
	(
		azure network public-ip delete --name $name --resource-group $group
	)
}
