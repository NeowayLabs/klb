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
