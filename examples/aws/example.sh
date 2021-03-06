#!/usr/bin/env nash

import nashlib/all
import klb/aws/all

defTags = (
	(Name klb-example)
	(Env testing)
)

fn print_resource(name, id) {
	printf "Created %s: %s%s%s\n" $name $NASH_GREEN $id $NASH_RESET
}

fn create_prod() {
	vpcid  <= aws_vpc_create("10.0.0.1/16", $defTags)
	appnet <= aws_subnet_create("10.0.1.0/24", $vpcid, $defTags)
	dbnet  <= aws_subnet_create("10.0.2.0/24", $vpcid, $defTags)
	igwid  <= aws_igw_create($defTags)
	tblid  <= aws_routetbl_create($vpcid, $defTags)

	aws_igw_attach($igwid, $vpcid)
	aws_route2igw($tblid, "0.0.0.0/0", $igwid)

	grpid <= aws_secgroup_create("klb-default-sg", "sg description", $vpcid, $defTags)

	print_resource("VPC", $vpcid)
	print_resource("app subnet", $appnet)
	print_resource("db subnet", $dbnet)
	print_resource("Internet Gateway", $igwid)
	print_resource("Routing table", $tblid)
	print_resource("Security group", $grpid)
}

create_prod()
