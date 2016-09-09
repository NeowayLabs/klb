#!/usr/bin/env nash

vpcTags = (
	(Name klb-vpc-example)
	(Env testing)
)

igwTags = (
	(Name klb-igw-example)
	(Env testing)
)

routeTblTags = (
	(Name klb-rtbl-example)
        (ENV testing)
)

appSubnetTags = (
	(Name klb-app-subnet-example)
        (Env testing)
)

dbSubnetTags = (
	(Name klb-db-subnet-example)
        (Env testing)
)

fn create_prod() {
	vpcid	<= aws_vpc_create("10.0.0.1/16", $vpcTags)
	appnet	<= aws_subnet_create("10.0.1.0/24", $appSubnetTags)
	dbnet	<= aws_subnet_create("10.0.2.0/24", $dbSubnetTags)
	igwid	<= aws_igw_create($igwTags)
	tblid	<= aws_routetbl_create($vpcid, $routeTblTags)

	aws_igw_attach($igwid, $vpcid)
        aws_route2igw(tblid, "0.0.0.0/0", $igwid)

	machine-up ./machines/app.yml
	machine-up ./machines/db.yml

	# TODO
}

create_prod()
