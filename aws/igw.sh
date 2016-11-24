# Internet gateway related functions

fn aws_igw_create(tags) {
	igwId <= (
		aws ec2 create-internet-gateway |
		jq ".InternetGateway.InternetGatewayId" |
		xargs echo -n
	)

	aws_tag($igwId, $tags)

	return $igwId
}

fn aws_igw_delete(igwId) {
	aws ec2 delete-internet-gateway --internet-gateway $igwId
}

fn aws_igw_attach(igwId, vpcId) {
	aws ec2 attach-internet-gateway --internet-gateway-id $igwId --vpc-id $vpcId >[1=]
}

fn aws_igw_detach(igwId, vpcId) {
	aws ec2 detach-internet-gateway --internet-gateway-id $igwId --vpc-id $vpcId
}

fn aws_igw_info(igwId) {
	info <= aws ec2 describe-internet-gateways

	return $info
}
