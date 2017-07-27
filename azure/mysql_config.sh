# Mysql Servers Configuration related functions

# azure_mysql_config_set sets a new configuration for a Mysql Server.
# On failure it will return en error string with details
# `name` is the name of the configuration
# `group` is the name of resource group.
# `servername` is the name of the Mysql server.
# `value` is the value of configuration
fn azure_mysql_config_set(name, group, servername, value) {
	_, status <= (
		az mysql server configuration set 
						--name $name
						--resource-group $group
						--server-name $servername
						--value $value
	)

	if $status != "0" {
		return format("unable to set mysql configuration[%s]", $name)
	}
}
