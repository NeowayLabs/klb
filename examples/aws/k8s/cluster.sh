#!/usr/bin/env nash

# Based on documentation of Kelsey Hightower - Kubernetes The Hard Way
# https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/01-infrastructure-aws.md

import klb/aws/all

k8sTags = (Name klb-kubernetes)

fn create_network() {
	vpcid    <= aws_vpc_create("10.240.0.0/16", $k8sTags)
	dhcpid   <= aws_dhcp_createopt("us-west-2.compute.internal", "AmazonProvidedDNS", $k8sTags)
	subnetid <= aws_subnet_create("10.240.0.0/24", $vpcid, $k8sTags)
	igwid    <= aws_igw_create($k8sTags)
	rtblid   <= aws_routetbl_create($vpcid, $k8sTags)

	aws_vpc_enabledns($vpcid, "enable-hostnames")
	aws_dhcp_assoc($dhcpid, $vpcid)
	aws_igw_attach($igwid, $vpcid)
	aws_routetbl_assoc($rtblid, $subnetid)
	aws_route2igw($rtblid, "0.0.0.0/0", $igwid)

	# Firewall rules
	secgrpName = "kubernetes"
	secgrpDesc = "Kubernetes security group"

	secgrpid   <= aws_secgroup_create($secgrpName, $secgrpDesc, $vpcid, $k8sTags)
}
