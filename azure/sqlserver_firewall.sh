# SQL Servers firewall related functions

# azure_sqlserver_firewall_rule_create creates a new rule
# for the firewall of a SQL Server.
# On failure it will return en error string with details
# `name` is the name of the rule.
# `group` is the name of resource group.
# `servername` is the name of the SQL Server.
# `startip` is the start IP address of the firewall rule. (IPV4 format)
# `endip` is the end IP address of the firewall rule. (IPV4 format)
fn azure_sqlserver_firewall_rule_create(name, group, servername, startip, endip) {
	_, status <= (
		az sql server firewall-rule create
						--name $name
						--resource-group $group
						--server $servername
						--start-ip-address $startip
						--end-ip-address $endip
	)

	if $status != "0" {
		return format("unable to create firewall rule[%s]", $name)
	}
}

# azure_sqlserver_firewall_rule_update updates a rule
# for the firewall of a SQL Server.
# On failure it will return en error string with details
# `name` is the name of the rule.
# `group` is the name of resource group.
# `servername` is the name of the SQL Server.
# `startip` is the start IP address of the firewall rule. (IPV4 format)
# `endip` is the end IP address of the firewall rule. (IPV4 format)
fn azure_sqlserver_firewall_rule_update(name, group, servername, startip, endip) {
	_, status <= (
		az sql server firewall-rule update
						--name $name
						--resource-group $group
						--server $servername
						--start-ip-address $startip
						--end-ip-address $endip
	)

	if $status != "0" {
		return format("unable to update firewall rule[%s]", $name)
	}
}

# azure_sqlserver_firewall_rule_delete deletes
# for the firewall of a SQL Server.
# On failure it will return en error string with details
# `name` is the name of the rule.
# `group` is the name of resource group.
# `servername` is the name of the SQL Server.
fn azure_sqlserver_firewall_rule_delete(name, group, servername) {
	_, status <= (
		az sqlserver server firewall-rule delete
						-y
						--name
						$name --resource-group $group --server $servername
	)

	if $status != "0" {
		return format("unable to delete firewall rule[%s]", $name)
	}
}
