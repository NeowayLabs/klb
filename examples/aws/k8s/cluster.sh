#!/usr/bin/env nash

# Based on documentation of Kelsey Hightower - Kubernetes The Hard Way
# https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/01-infrastructure-aws.md
#
# Dependencies: nash, nashlib, awscli, jq, openssl

import klb/aws/all

# Key pair name
keyName = "kubernetes"

# Ubuntu Amazon EC2
imageid = "ami-746aba14"

# Default AWS tag used in all resources created
defTag = (environment KLB-kubernetes)

# AWS Tag used in network devices
netTags = (
	(Name "KLB-network")
	$defTag
)

fn clone(list) {
	other = ()

	for l in $list {
		other <= append($other, $l)
	}

	return $other
}

fn start_ssh_agent() {
	sshenv <= ssh-agent -s

	(
		echo $sshenv |
		sed "s/export/setenv/g" |
		sed "s/=/=\"/g" |
		sed "s/SSH_AUTH_SOCK;/SSH_AUTH_SOCK\\n/g" |
		sed "s/SSH_AGENT_PID;/SSH_AGENT_PID\\n/g" |
		sed "s/echo.*//g" |
		sed "s/;/\"\\n/g"
				> /tmp/ssh.env
	)

	cat /tmp/ssh.env
}

fn extract_public_key(pubKeyPath) {
	pubKey <= cat $pubKeyPath | grep -v " KEY" | xargs echo -n | sed "s/ //g"

	return $pubKey
}

fn setup_ssh() {
	privKeyPath = $HOME+"/.ssh/"+$keyName
	pubKeyPath  = $privKeyPath+".pub"

	-test -f $keyPairPath

	if $status != "0" {
		aws_keypair_delete($keyName)

		openssl genrsa -out $privKeyPath 2048
		openssl rsa -in $privKeyPath -pubout > $pubKeyPath

		pubKey  <= extract_public_key($pubKeyPath)
		fingerp <= aws_keypair_import($keyName, $pubKey)

		printf "Pubkey %s imported\n" $fingerp
		chmod 600 $privKeyPath
	}

	# Check env var exists
	-echo $SSH_AUTH_SOCK >[1=]

	if $status != "0" {
		start_ssh_agent()

		# load agent environment vars
		import "/tmp/ssh.env"
	}

	ssh-add $privKeyPath
}

fn create_network() {
	vpcid    <= aws_vpc_create("10.240.0.0/16", $netTags)
	dhcpid   <= aws_dhcp_createopt("us-west-2.compute.internal", "AmazonProvidedDNS", $netTags)
	subnetid <= aws_subnet_create("10.240.0.0/24", $vpcid, $netTags)
	igwid    <= aws_igw_create($netTags)
	rtblid   <= aws_routetbl_create($vpcid, $netTags)

	aws_vpc_enabledns($vpcid, "enable-hostnames")
	aws_dhcp_assoc($dhcpid, $vpcid)
	aws_igw_attach($igwid, $vpcid)
	aws_routetbl_assoc($rtblid, $subnetid)
	aws_route2igw($rtblid, "0.0.0.0/0", $igwid)

	# Firewall rules
	secgrpName = "kubernetes"
	secgrpDesc = "Kubernetes security group"

	secgrpid   <= aws_secgroup_create($secgrpName, $secgrpDesc, $vpcid, $netTags)

	aws_secgroup_ingress($secgrpid, "all", "0-65535", "10.240.0.0/16")
	aws_secgroup_ingress($secgrpid, "tcp", "22", "0.0.0.0/0")
	aws_secgroup_ingress($secgrpid, "tcp", "6443", "0.0.0.0/0")

	# Kubernetes Public Address
	elbPort   = "6443"
	elbProto  = "TCP"
	listeners = "Protocol="+$elbProto+",LoadBalancerPort="+$elbPort
	listeners = $listeners+",InstanceProtocol="+$elbProto
	listeners = $listeners+",InstancePort="+$elbPort

	dnsName   <= aws_elb_create("kubernetes", $listeners, $subnetid, "internal", $secgrpid)

	printf "Kubernetes Public DNS: %s\n" $dnsName
	printf "Kubernetes network created successfully\n"

	setenv secgrpid
	setenv subnetid
	setenv vpcid
	setenv dhcpid
	setenv igwid
	setenv rtblid
	setenv dnsName
}

fn create_iam_policies() {
	roleid    <= aws_iam_create("kubernetes", "kubernetes-iam-role.json")
	policyArn <= aws_iam_createpolicy("kubernetes", "kubernetes-iam-policy.json")

	aws_iam_attachpolicy("kubernetes", $policyArn)
	aws_iam_profile("kubernetes")
	aws_iam_addrole2profile("kubernetes", "kubernetes")
}

fn create_etcdcluster() {
	etcd0Tags = (
		(Name "KLB-etcd0")
		$defTag
	)

	etcd1Tags = (
		(Name "KLB-etcd1")
		$defTag
	)

	etcd2Tags = (
		(Name "KLB-etcd2")
		$defTag
	)

	# Base configuration of all etcd instances
	etcdbase <= aws_instance_new($imageid, $keyName, "t2.small")
	etcdbase <= aws_instance_setsecgrps($etcdbase, $secgrpid)
	etcdbase <= aws_instance_setpubip($etcdbase)
	etcdbase <= aws_instance_setsubnet($etcdbase, $subnetid)
	etcdbase <= aws_instance_setcount($etcdbase, "1")
	etcd0    <= clone($etcdbase)
	etcd1    <= clone($etcdbase)
	etcd2    <= clone($etcdbase)

	# Configure network address of each etcd
	etcd0 <= aws_instance_setprivip($etcd0, "10.240.0.10")
	etcd1 <= aws_instance_setprivip($etcd1, "10.240.0.11")
	etcd2 <= aws_instance_setprivip($etcd2, "10.240.0.12")

	# Run ETCD cluster
	etcd0id <= aws_instance_run($etcd0, $etcd0Tags)
	etcd1id <= aws_instance_run($etcd1, $etcd1Tags)
	etcd2id <= aws_instance_run($etcd2, $etcd2Tags)
}

fn create_controllers() {
	# Kubernetes controllers
	ctl0Tags = (
		(Name "KLB-controller0")
		$defTag
	)

	ctl1Tags = (
		(Name "KLB-controller1")
		$defTag
	)

	ctl2Tags = (
		(Name "KLB-controller2")
		$defTag
	)

	# Setup base controller instance
	ctlbase <= aws_instance_new($imageid, $keyName, "t2.small")
	ctlbase <= aws_instance_setpubip($ctlbase)
	ctlbase <= aws_instance_setcount($ctlbase, "1")
	ctlbase <= aws_instance_setprofile($ctlbase, "Name=kubernetes")
	ctlbase <= aws_instance_setsecgrps($ctlbase, $secgrpid)
	ctlbase <= aws_instance_setsubnet($ctlbase, $subnetid)
	ctl0    <= clone($ctlbase)
	ctl1    <= clone($ctlbase)
	ctl2    <= clone($ctlbase)

	# Set specific controller config
	ctl0 <= aws_instance_setprivip($ctl0, "10.240.0.20")
	ctl1 <= aws_instance_setprivip($ctl1, "10.240.0.21")
	ctl2 <= aws_instance_setprivip($ctl2, "10.240.0.22")

	# Start controllers
	ctl0id <= aws_instance_run($ctl0, $ctl0Tags)
	ctl1id <= aws_instance_run($ctl1, $ctl1Tags)
	ctl2id <= aws_instance_run($ctl2, $ctl2Tags)

	# Modify attributes
	aws_instance_modify($ctl0id, "sourceDestCheck", "false")
	aws_instance_modify($ctl1id, "sourceDestCheck", "false")
	aws_instance_modify($ctl2id, "sourceDestCheck", "false")
}

fn create_workers() {
	worker0Tags = (
		(Name "KLB-worker0")
		$defTag
	)

	worker1Tags = (
		(Name "KLB-worker1")
		$defTag
	)

	worker2Tags = (
		(Name "KLB-worker2")
		$defTag
	)

	# base worker setup
	workerbase <= aws_instance_new($imageid, $keyName, "t2.small")
	workerbase <= aws_instance_setpubip($workerbase)
	workerbase <= aws_instance_setprofile($workerbase, "Name=kubernetes")
	workerbase <= aws_instance_setcount($workerbase, "1")
	workerbase <= aws_instance_setsecgrps($workerbase, $secgrpid)
	workerbase <= aws_instance_setsubnet($workerbase, $subnetid)

	# workers
	worker0 <= clone($workerbase)
	worker1 <= clone($workerbase)
	worker2 <= clone($workerbase)
	worker0 <= aws_instance_setprivip($worker0, "10.240.0.30")
	worker1 <= aws_instance_setprivip($worker1, "10.240.0.31")
	worker2 <= aws_instance_setprivip($worker2, "10.240.0.32")

	# Run
	worker0id <= aws_instance_run($worker0, $worker0Tags)
	worker1id <= aws_instance_run($worker1, $worker1Tags)
	worker2id <= aws_instance_run($worker2, $worker2Tags)

	# modify
	aws_instance_modify($worker0id, "sourceDestCheck", "false")
	aws_instance_modify($worker1id, "sourceDestCheck", "false")
	aws_instance_modify($worker2id, "sourceDestCheck", "false")
}

fn create_instances() {
	create_etcdcluster()
	create_controllers()
	create_workers()

	filters = (
		("tag:environment" "KLB-kubernetes")
	)

	instances <= aws_instance_describe($filters)

	for instance in $instances {
		printf "%s%s%s\n" $NASH_GREEN $instance $NASH_RESET
	}
}

setup_ssh()
create_network()
create_iam_policies()
create_instances()
