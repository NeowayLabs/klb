#!/usr/bin/env nash

import klb/aws

igwTags = (
    ("Name" "igw-core")
    ("Env" "Production")
)

vpcTags = (
    ("Name" "vpc-core")
    ("Env" "Production")
)

# Create an Internet Gateway and attach a vpc to it
vpcid <= vpc_create("10.0.0.0/16", $vpcTags)
igwid <= igw_create($igwTags)
subid <= subnet_create($vpcid, "10.0.1.0/24")

attach_igw($igwid, $vpcid)

machine-up ./machines/node1.yml
machine-up ./machines/node2.yml
machine-up ./machines/node3.yml
