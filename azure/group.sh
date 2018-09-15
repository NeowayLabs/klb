# Resource Group related functions

# azure_group_create creates a new `resource group`.
# `name` is the resource group name
# `location` is the azure region
fn azure_group_create(name, location) {
	azure group create --name $name --location $location
}

# azure_group_delete deletes a `resource group`.
# `name` is the resource group name
fn azure_group_delete(name) {
	azure group delete -q --name $name
}

# azure_group_delete_async deletes a `resource group`.
# `name` is the resource group name
#
# This function will return imediately (wont wait the delete operation to finish).
fn azure_group_delete_async(name) {
	out, status <= az group delete --yes --no-wait --name $name
    if $status != "0" {
        format("error deleting resgroup[%s]: %s", $name, $out)
    }
    return ""
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

# azure_group_location gets the location of a given resgroup.
#
# It returns two values, the location and an empty error string on success or
# an empty string and an error string if something goes wrong.
fn azure_group_location(name) {
    out, status <= az group show --name $name >[2=1]
    if $status != "0" {
        return "", format("error loading resgroup[%s] location: %s", $name, $out)
    }

    location <= echo $out | jq -r ".location"
    return $location, ""
}
