# Group related functions

fn azure_group_create(name, location) {
	azure group create --name $name --location $location
}

fn azure_group_delete(name) {
	azure group delete -q --name $name
}
