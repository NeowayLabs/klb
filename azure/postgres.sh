# Postgres Servers related functions

# azure_postgres_new creates a new server of "managed Postgres".
# `name` is the name of the server.
# `group` is name of resource group.
# `location` is the Azure Region.
# `username` is the login name of a server.
# `password` is the postgres server login password.
fn azure_postgres_new(name, group, location, username, password) {
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

# azure_postgres_set_compute_units sets the number of compute units.
# `instance` is the postgres server instance.
# `units` is the number of compute units desired.
# Default: 100
fn azure_postgres_set_compute_units(instance, units) {
	instance <= append($instance, "--compute-units")
	instance <= append($instance, $units)

	return $instance
}

# azure_postgres_set_performance_tier sets the performance tier.
# `instance` is the postgres server instance.
# `tier` is the performance tier desired for the server.
# Allowed values: Basic, Standard.
# Default: Basic
fn azure_postgres_set_performance_tier(instance, tier) {
	instance <= append($instance, "--performance-tier")
	instance <= append($instance, $tier)

	return $instance
}

# azure_postgres_disable_ssl disables the ssl enforcement.
# `instance` is the postgres server instance.
fn azure_postgres_disable_ssl(instance) {
	instance <= append($instance, "--ssl")
	instance <= append($instance, "Disabled")

	return $instance
}

# azure_postgres_set_max_size sets the max storage size for the server.
# `instance` is the postgres server instance.
# `size` is the postgres server storage size. (in MB)
# Default: 51200
fn azure_postgres_set_max_size(instance, size) {
	instance <= append($instance, "--storage-size")
	instance <= append($instance, $size)

	return $instance
}

# azure_postgres_set_version sets the Postgres version of a server.
# `instance` is the postgres server instance.
# `version` is the postgres server version.
# Default: 9.5
fn azure_postgres_set_version(instance, version) {
	instance <= append($instance, "--version")
	instance <= append($instance, $version)

	return $instance
}

# azure_postgres_server_create creates a new "managed Postgres server"
# `instance` is the postgres server instance.
fn azure_postgres_server_create(instance) {
        az postgres server create $instance
}
