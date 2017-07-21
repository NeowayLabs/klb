# Mysql Servers related functions

# azure_mysql_new creates a new server of "managed Myqsl".
# `name` is the name of the server.
# `group` is name of resource group.
# `location` is the Azure Region.
# `username` is the login name of a server.
# `password` is the mysql server login password.
fn azure_mysql_new(name, group, location, username, password) {
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

# azure_mysql_set_compute_units sets the number of compute units.
# `instance` is the mysql server instance.
# `units` is the number of compute units desired.
# Default: 100
fn azure_mysql_set_compute_units(instance, units) {
	instance <= append($instance, "--compute-units")
	instance <= append($instance, $units)

	return $instance
}

# azure_mysql_set_performance_tier sets the performance tier.
# `instance` is the mysql server instance.
# `tier` is the performance tier desired for the server.
# Allowed values: Basic, Standard.
# Default: Basic
fn azure_mysql_set_performance_tier(instance, tier) {
	instance <= append($instance, "--performance-tier")
	instance <= append($instance, $tier)

	return $instance
}

# azure_mysql_disable_ssl disables the ssl enforcement.
# `instance` is the mysql server instance.
fn azure_mysql_disable_ssl(instance) {
	instance <= append($instance, "--ssl-enforcement")
	instance <= append($instance, "Disabled")

	return $instance
}

# azure_mysql_set_storage_size sets the max storage size for the server.
# `instance` is the mysql server instance.
# `size` is the mysql server storage size. (in MB)
# Default: 50GB
fn azure_mysql_set_storage_size(instance, size) {
	instance <= append($instance, "--storage-size")
	instance <= append($instance, $size)

	return $instance
}

# azure_mysql_set_version sets the Mysql version of a server.
# `instance` is the mysql server instance.
# `version` is the mysql server version.
# Default: 5.7
fn azure_mysql_set_version(instance, version) {
	instance <= append($instance, "--version")
	instance <= append($instance, $version)

	return $instance
}

# azure_mysql_set_tags sets the tags for the server.
# `instance` is the mysql server instance.
# `tags` is the tags of mysql server.
fn azure_mysql_set_storage_size(instance, tags) {
	instance <= append($instance, "--tags")
	instance <= append($instance, $tags)

	return $instance
}

# azure_mysql_server_create creates a new "managed Mysql server"
# `instance` is the mysql server instance.
fn azure_mysql_server_create(instance) {
	az mysql server create $instance
}
