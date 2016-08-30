#!/usr/bin/env nash

IFS = ()

# Create a VPC
fn vpc_create(cidr) {
	vpcid <= aws ec2 create-vpc --cidr-block $cidr | jq ".Vpc.VpcId"
	return $vpcid
}

# Create an Internet Gateway
fn igw_create() {
	igwid <= aws ec2 create-internet-gateway | jq ".InternetGateway.InternetGatewayId"
	return $igwid
}

# Attach an internet gateway to a VPC
fn attach_igw(igwid, vpcid) {
	# do not generate output
	aws ec2 attach-internet-gateway --internet-gateway-id $igwid --vpc-id $vpcid
}

# Create a subnet
fn subnet_create(vpcid, cidr) {
	netid <= aws ec2 create-subnet --vpc-id $vpcid --cidr-block $cidr | jq ".Subnet.SubnetId"
	return $netid
}
