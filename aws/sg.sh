# security groups related functions

# Create a security group on EC2Classic or EC2VPC
# vpcid is not empty.
fn aws_secgroup_create(name, desc, vpcid, tags) {
	vpcarg = ()

	if $vpcid != "" {
		vpcarg = ("--vpc-id" $vpcid)
	}

	grpid <= (
		aws ec2 create-security-group
					--group-name $name
					--description $desc $vpcarg |
		jq ".GroupId" |
		xargs echo -n
	)

	aws_tag($grpid, $tags)

	return $grpid
}

fn aws_secgroup_delete(grpid) {
	aws ec2 delete-security-group --group-id $grpid
}

fn aws_secgroup_ingress(grpid, proto, port, cidr) {
	(
		aws ec2 authorize-security-group-ingress
							--group-id $grpid
							--protocol $proto
							--port $port
							--cidr $cidr
							>[1=]
	)
}
