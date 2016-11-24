# Subnet related functions

fn aws_subnet_create(cidr, vpcid, tags) {
	netid <= (
		aws ec2 create-subnet
				--vpc-id $vpcid
				--cidr-block $cidr |
		jq ".Subnet.SubnetId" |
		xargs echo -n
	)

	aws_tag($netid, $tags)

	return $netid
}

fn aws_subnet_delete(netid) {
	aws ec2 delete-subnet --subnet $netid >[1=]
}
