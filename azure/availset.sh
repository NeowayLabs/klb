# Availset related functions

# azure_availset_create creates a new `availability set`.
# `name` is the availability set name
# `group` is the resource group name
# `location` is the azure region
fn azure_availset_create(name, group, location) {
	azure availset create --name $name --resource-group $group --location $location
}

# azure_availset_delete deletes a exit `availability set`.
# `name` is the availability set name
# `group` is the resource group name
fn azure_availset_delete(name, group) {
	azure availset delete -q --name $name --resource-group $group
}

# azure_availset_exists checks if a availability set exists
# `name` is the availability set name
# `group` is the resource group name
#
# If the availability set exists returns "0", otherwise returns
# a value that is not "0" (like process status codes).
fn azure_availset_exists(name, group) {
	out, status <= az vm availability-set show --name $name --resource-group $group
	if $status != "0" {
		print("error[%s]: avset[%s] resgroup[%s] do not exist\n", $out, $name, $group)
		return $status
	}
	# WHY: az succeeds but the stdout is empty when it does not exists lol
	if $out == "" {
		return "1"
	}
	return "0"
}
