# Machine/instances related functions

fn aws_instance_create(imgid, kname, secgrps, type, privip, subnetid, tags) {
	instid <= (
		aws ec2 run-instances
				--image-id $imgid
				--count 1
				--key-name kubernetes
				--security-group-ids $secgrps
				--instance-type $type
				--private-ip-address $privip
				--subnet-id $subnetid
				--associate-public-ip-address
				 |
		jq -r ".Instances[].InstanceId" | xargs echo -n
	)

	aws_tag($instid, $tags)

	return $instid
}
