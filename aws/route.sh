# Routing table specific functions

fn aws_routetbl_create(vpcid, tags) {
	tblid <= (
		aws ec2 create-route-table
					--vpc-id $vpcid |
		jq ".RouteTable.RouteTableId" |
		xargs echo -n
	)

	aws_tag($tblid, $tags)

	return $tblid
}

fn aws_route2igw(tblid, cidr, igwid) {
	(
		aws ec2 create-route
				--route-table-id $tblid
				--destination-cidr-block $cidr
				--gateway-id $igwid
				>[1=]
	)
}

fn aws_route2vpcpeer(tblid, cidr, peerid) {
	(
		aws ec2 create-route
				--route-table-id $tblid
				--destination-cidr-block $cidr
				--vpc-peering-connection-id $peerid
	)
}

fn aws_route2nat(tblid, cidr, natid) {
	(
		aws ec2 create-route
				--route-table-id $tblid
				--destination-cidr-block $cidr
				--nat-gateway-id $natid
	)
}

fn aws_route2netif(tblid, cidr, netifid) {
	(
		aws ec2 create-route
				--route-table-id $tblid
				--destination-cidr-block $cidr
				--network-interface-id $netifid
	)
}

fn aws_route2instance(tblid, cidr, instid) {
	(
		aws ec2 create-route
				--route-table-id $tblid
				--destination-cidr-block $cidr
				--instance-id $instid
	)
}

fn aws_routetbl_assoc(rtblid, netid) {
	aws ec2 associate-route-table --route-table-id $rtblid --subnet-id $netid >[1=]
}
