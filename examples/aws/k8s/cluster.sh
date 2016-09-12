#!/usr/bin/env nash

import klb/aws/all

vpcTags = ((Name kubernetes))

fn create_network() {
	vpcid <= aws_vpc_create("10.240.0.0/16", $vpcTags)

	aws_vpc_enabledns($vpcid, "enable-hostnames")
}
