# Availset related functions

fn azure_availset_create(name, group, location) {
	azure availset create 
		--name $name
		--resource-group $group 
		--location $location
}

fn azure_availset_delete(name) {
	azure availset delete 
		--name $name
		--resource-group $group
}