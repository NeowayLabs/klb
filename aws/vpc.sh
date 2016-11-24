# VPC related functions

# Creates a vpc setting its name, cidr and tags
fn aws_vpc_create(cidr, tags) {
	vpcid <= (
		aws ec2 create-vpc
				--cidr-block $cidr |
		jq ".Vpc.VpcId" |
		xargs echo -n
	)

	aws_tag($vpcid, $tags)

	return $vpcid
}

# Delete a vpc
fn aws_vpc_delete(vpcid) {
	aws ec2 delete-vpc --vpc-id $vpcid
}

fn aws_vpc_info(vpcid) {
	json <= -aws ec2 describe-vpcs --vpc-id $vpcid >[2=]

	return $json
}

fn aws_vpc_enabledns(vpcid, enableHostname) {
	paramStr = "{\"Value\": true}"

	aws ec2 modify-vpc-attribute --vpc-id $vpcid --enable-dns-support $paramStr

	if $enableHostname != "" {
		aws ec2 modify-vpc-attribute --vpc-id $vpcid --enable-dns-hostnames $paramStr
	}
}
