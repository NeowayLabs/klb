# Active Directory applications related functions

# azure_app_create creates a new `active directory application`.
# `name` is the display name for the application
# `homepage` is the URL to the application homepage
# `uri` is the comma-delimitied URIs that identify the application
fn azure_app_create(name, homepage, uri) {
	out   <= (
		azure ad app create --json
					--name $name
					--home-page $homepage
					--identifier-uris $uri
	)

	appid <= echo -n $out | jq -r ".appId" | tr -d "\n"
	objid <= echo -n $out | jq -r ".objectId" | tr -d "\n"

	ret   = ($appid $objid)

	return $ret
}

# azure_app_delete deletes a exit `active directory application`.
# `objid` is the display name for the application
fn azure_app_delete(objid) {
	azure ad app delete --objectid $objid
}
