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

fn azure_app_delete(objid) {
	azure ad app delete --objectid $objid
}
