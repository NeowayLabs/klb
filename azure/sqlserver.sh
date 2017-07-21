# SQL Servers related functions

# azure_sqlserver_new creates a new Logical SQL Server on Azure.
# `name` is the name of the server.
# `group` is name of resource group.
# `location` is the Azure Region.
# `username` is the login name of a server.
# `password` is the postgres server login password.
fn azure_sqlserver_new(name, group, location, username, password) {
	instance = (
		"--name"
		$name
		"--resource-group"
		$group
		"--location"
		$location
		"--admin-user"
		$username
		"--admin-password"
		$password
	)

	return $instance
}

# azure_sqlserver_set_tags set the resource tags
# `instance` is the sqlserver instance.
# `tag` is the tag desired.
fn azure_sqlserver_set_tags(instance, tags) {
	instance <= append($instance, "--tags")
	instance <= append($instance, $tags)

	return $instance
}

# azure_sqlserver_server_create creates a new "managed SQLServer Logical Server"
# `instance` is the sqlserver server instance.
fn azure_sqlserver_server_create(instance) {
	az sql server create $instance
}
