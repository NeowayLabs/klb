#!/usr/bin/env nash

import klb/aws/all

igwTags = ((Name klb-igw-tests) (Env TEST))
vpcTags = ((Name klb-vpc-igw-tests) (Env TEST))

fn create() {
	igwId <= aws_igw_create($igwTags)

	echo "Internet gateway created: " $igwId

	return $igwId
}

fn test_igw() {
	error = "0"

	igwId <= create()
	info <= aws_igw_info($igwId)
	info <= echo -n $info | jq ".InternetGateways[0].Attachments"

	echo $info | -grep "VpcId"

	if $status == "0" {
		echo "IGW created with attachments... :: " $info

		aws_igw_delete($igwId)

		abort
	}

	vpcId <= aws_vpc_create("10.0.0.1/16", $vpcTags)

	aws_igw_attach($igwId, $vpcId)

	vpcAttached <= aws_igw_info($igwId)
	vpcAttached <= (
		echo -n $vpcAttached |
		jq ".InternetGateways[0].Attachments[0].VpcId" |
		xargs echo -n
	)

	if $status != "0" {
		printf "Failed to attach igw %s to vpc %s" $igwId $vpcId

		error = "1"
	}

	# this order matters
	aws_igw_detach($igwId, $vpcId)
	aws_vpc_delete($vpcId)
	aws_igw_delete($igwId)

	if $error != "0" {
		abort
	}
}

test_igw()
