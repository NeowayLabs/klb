# VPC related functions

# Creates a vpc setting its name, cidr and tags
fn aws_vpc_create(cidr, tags) {
	vpcId <= (
		aws ec2 create-vpc --cidr-block $cidr |
		jq ".Vpc.VpcId" |
		xargs echo -n
	)

	aws_tag($vpcId, $tags)

        return $vpcId
}

# Delete a vpc
fn aws_vpc_delete(vpcId) {
	aws ec2 delete-vpc --vpc-id $vpcId
}

fn aws_vpc_info(vpcId) {
	IFS = ()
   	json <= -aws ec2 describe-vpcs --vpc-id $vpcId >[2=]

        return $json
}
