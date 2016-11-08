# Network Security Group related functions

fn azure_nsg_create(name, group, location) {
	azure network nsg create 
		--name $name
		--resource-group $group
		--location $location
}

fn azure_nsg_delete(name, group, location) {
	azure network nsg delete 
		--name $name
		--resource-group $group
}

fn azure_nsg_add_rule(name, group, nsgname, direction, priority, protocol, address, port, access) {
	azure network nsg rule create 
	    --name $name
		--resource-group $group
		--nsg-name $nsgname
	    --direction $direction
	    --priority $priority
	    --protocol $protocol
	    --source-address-prefix $address
	    --destination-port-range $port
	    --access $access
}

fn azure_nsg_delete_rule(name, group, nsgname) {
	azure network nsg rule delete 
	    --name $name
		--resource-group $group
		--nsg-name $nsgname
}