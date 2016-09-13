#!/usr/bin/env nash

# Based on documentation of Kelsey Hightower - Kubernetes The Hard Way
# https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/01-infrastructure-aws.md

import klb/aws/all

tags = (
	(Name klb-kubernetes)
)

fn create_network() {
	vpcid    <= aws_vpc_create("10.240.0.0/16", $tags)
	dhcpid   <= aws_dhcp_createopt("us-west-2.compute.internal", "AmazonProvidedDNS", $tags)
	subnetid <= aws_subnet_create("10.240.0.0/24", $vpcid, $tags)
	igwid    <= aws_igw_create($tags)
	rtblid   <= aws_routetbl_create($vpcid, $tags)

	aws_vpc_enabledns($vpcid, "enable-hostnames")
	aws_dhcp_assoc($dhcpid, $vpcid)
	aws_igw_attach($igwid, $vpcid)
	aws_routetbl_assoc($rtblid, $subnetid)
	aws_route2igw($rtblid, "0.0.0.0/0", $igwid)

	# Firewall rules
	secgrpName = "kubernetes"
	secgrpDesc = "Kubernetes security group"

	secgrpid   <= aws_secgroup_create($secgrpName, $secgrpDesc, $vpcid, $tags)

	aws_secgroup_ingress($secgrpid, "all", "0-65535", "10.240.0.0/16")
	aws_secgroup_ingress($secgrpid, "tcp", "22", "0.0.0.0/0")
	aws_secgroup_ingress($secgrpid, "tcp", "6443", "0.0.0.0/0")

	# Kubernetes Public Address
	elbPort   = "6443"
	elbProto  = "TCP"
	listeners = "Protocol="+$elbProto+",LoadBalancerPort="+$elbPort
	listeners = $listeners+",InstanceProtocol="+$elbProto
	listeners = $listeners+",InstancePort="+$elbPort

	dnsName   <= aws_elb_create("kubernetes", $listeners, $subnetid, $secgrpid)

	printf "Kubernetes Public DNS: %s\n" $dnsName
	printf "Kubernetes network created successfully\n"
}

create_network()

#create_instances()
