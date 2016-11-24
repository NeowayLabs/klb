# NAT gateway related functions

# Create a NAT gateway
fn aws_nat_create(subnetId, eipId, tags) {
	natId <= (
		aws ec2 create-nat-gateway
					--subnet-id $subnetId
					--allocation-id $eipId |
		jq ".NatGateway.NatGatewayId" |
		xargs echo -n
	)

	aws_tags($natId, $tags)

	return $natId
}
