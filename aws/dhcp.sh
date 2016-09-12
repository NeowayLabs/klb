# DHCP related functions

fn aws_dhcp_createopt(domain, domainServers, tags) {
	dhcpOptId <= (
		aws ec2 create-dhcp-options
					--dhcp-configuration "Key=domain-name,Values="+$domain "Key=domain-name-servers,Values="+$domainServers |
		jq ".DhcpOptions.DhcpOptionsId" |
		xargs echo -n
	)

	aws_tag($dhcpOptId, $tags)
}

fn aws_dhcp_assoc(optid, vpcid) {
	aws ec2 associate-dhcp-options --dhcp-options-id $opt --vpc-id $vpcid
}
