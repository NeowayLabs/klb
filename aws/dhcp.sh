# DHCP related functions

fn aws_dhcp_createopt(domain, domainServers, tags) {
	config = (
		"Key=domain-name,Values="+$domain
		"Key=domain-name-servers,Values="+$domainServers
	)

	OptId <= (
		aws ec2 create-dhcp-options
					--dhcp-configuration $config |
		jq ".DhcpOptions.DhcpOptionsId" |
		xargs echo -n
	)

	aws_tag($optId, $tags)

	return $optId
}

fn aws_dhcp_assoc(optid, vpcid) {
	aws ec2 associate-dhcp-options --dhcp-options-id $optid --vpc-id $vpcid
}
