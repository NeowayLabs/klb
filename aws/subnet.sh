# Subnet related functions

fn aws_subnet_create(vpcId, cidr, tags) {
	netId <= (
		aws ec2 create-subnet	--vpc-id $vpcId
					--cidr-block $cidr |
		jq ".Subnet.SubnetId"
	)

	return $netId
}

fn aws_subnet_delete(netId) {
	aws ec2 delete-subnet --subnet $netId
}
