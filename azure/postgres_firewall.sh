# Postgres Servers firewall related functions

# azure_postgres_firewall_rule_create creates a new rule
# for the firewall of a server of "managed Postgres".
# `name` is the name of the rule.
# `group` is the name of resource group.
# `servername` is the name of the Postgres server.
# `startip` is the start IP address of the firewall rule. (IPV4 format)
# `endip` is the end IP address of the firewall rule. (IPV4 format)
fn azure_postgres_firewall_rule_create(name, group, servername, startip, endip) {
	(
		az postgres server
		   firewall-rule create
		   --name $name
		   --resource-group $group
		   --server-name $servername
                   --start-ip-address $startip
                   --end-ip-address $endip
	)
}

# azure_postgres_firewall_rule_update updates a rule
# for the firewall of a server of "managed Postgres".
# `name` is the name of the rule.
# `group` is the name of resource group.
# `servername` is the name of the Postgres server.
# `startip` is the start IP address of the firewall rule. (IPV4 format)
# `endip` is the end IP address of the firewall rule. (IPV4 format)
fn azure_postgres_firewall_rule_update(name, group, servername, startip, endip) {
	(
		az postgres server
		   firewall-rule update
		   --name $name
		   --resource-group $group
                   --server-name $servername
                   --start-ip-address $startip
                   --end-ip-address $endip
	)
}

# azure_postgres_firewall_rule_delete deletes
# for the firewall of a server of "managed Postgres".
# `name` is the name of the rule.
# `group` is the name of resource group.
# `servername` is the name of the Postgres server.
fn azure_postgres_firewall_rule_delete(name, group, servername) {
	(
		az postgres server
		   firewall-rule delete -y
		   --name $name
		   --resource-group $group
                   --server-name $servername
	)
}
