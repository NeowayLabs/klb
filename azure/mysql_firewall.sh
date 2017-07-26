# Mysql Servers firewall related functions

# azure_mysql_firewall_rule_create creates a new rule
# for the firewall of a server of "managed Mysql".
# On failure it will return en error string with details
# `name` is the name of the rule.
# `group` is the name of resource group.
# `servername` is the name of the Mysql server.
# `startip` is the start IP address of the firewall rule. (IPV4 format)
# `endip` is the end IP address of the firewall rule. (IPV4 format)
fn azure_mysql_firewall_rule_create(name, group, servername, startip, endip) {
	_, status <= (
		az mysql server firewall-rule create
						--name $name
						--resource-group $group
						--server-name $servername
						--start-ip-address $startip
						--end-ip-address $endip
	)

	if $status != "0" {
		return format("unable to create firewall rule[%s]", $name)
	}
}

# azure_mysql_firewall_rule_update updates a rule
# for the firewall of a server of "managed Mysql".
# On failure it will return en error string with details
# `name` is the name of the rule.
# `group` is the name of resource group.
# `servername` is the name of the Mysql server.
# `startip` is the start IP address of the firewall rule. (IPV4 format)
# `endip` is the end IP address of the firewall rule. (IPV4 format)
fn azure_mysql_firewall_rule_update(name, group, servername, startip, endip) {
	_, status <= (
		az mysql server firewall-rule update
						--name $name
						--resource-group $group
						--server-name $servername
						--start-ip-address $startip
						--end-ip-address $endip
	)

	if $status != "0" {
		return format("unable to update firewall rule[%s]", $name)
	}
}

# azure_mysql_firewall_rule_delete deletes
# for the firewall of a server of "managed Mysql".
# On failure it will return en error string with details
# `name` is the name of the rule.
# `group` is the name of resource group.
# `servername` is the name of the Mysql server.
fn azure_mysql_firewall_rule_delete(name, group, servername) {
	_, status <= (
		az mysql server firewall-rule delete
						-y
						--name
						$name --resource-group $group --server-name $servername
	)

	if $status != "0" {
		return format("unable to delete firewall rule[%s]", $name)
	}
}
