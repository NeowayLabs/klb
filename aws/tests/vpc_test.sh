#!/usr/bin/env nash

import klb/aws/all

vpcTags = (
	(Name klb-test-vpc)
	(Env TEST)
)

fn create() {
	vpcId <= aws_vpc_create("10.0.0.1/16", $vpcTags)

	echo "VPC ID created: " $vpcId
	return $vpcId
}

fn destroy(vpcId) {
	aws_vpc_delete($vpcId)
}

fn test_vpc() {
	vpcId   <= create()
	vpcInfo <= aws_vpc_info($vpcId)

	if $vpcInfo == "" {
		echo "Failed to create vpc"
		abort
	}

	IFS = ("\n")
	tagKeys <= echo $vpcInfo | jq ".Vpcs[].Tags[].Key"
        keyslen <= len($tagKeys)

	echo "Tag keys: " $tagKeys

	if $keyslen != "2" {
		echo "Failed to tag the vpc"
                echo "Found tags: " $tagKeys
                destroy($vpcId)
                return
	}

        tagValues <= echo $vpcInfo | jq ".Vpcs[].Tags[].Value"
        valueslen <= len($tagValues)

        if $valueslen != "2" {
		echo "Failed to tag the VPC"
                echo "Found tag values: " $tagValues
		destroy($vpcId)
                return
	}

	destroy($vpcId)

	vpcInfo <= aws_vpc_info($vpcId)

	if $vpcInfo != "" {
		echo "Failed to destroy vpc: " $vpcId
		abort
	}
}

test_vpc()
