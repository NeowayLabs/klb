# Resource Group related functions

# azure_group_create creates a new `resource group`.
# `name` is the resource group name
# `location` is the azure region
fn azure_group_create(name, location) {
	azure group create --name $name --location $location
}

# azure_group_delete deletes a exit `resource group`.
# `name` is the resource group name
fn azure_group_delete(name) {
	azure group delete -q --name $name
}

# azure_group_get_names returns a list with the names
# of all available resource groups
fn azure_group_get_names() {
	res <= az group list --query "[*].name" | jq -r ".[]"

	return split($res, "\n")
}

# azure_group_exists returns "0" if a resource group
# already exists, "1" otherwise.
fn azure_group_exists(name) {
	group_name <= az group show --name $name | jq -r ".name"

	if $group_name == $name {
		return "0"
	}

	return "1"
}
