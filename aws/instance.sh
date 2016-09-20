# Machine/instances related functions

fn aws_instance_new(imgid, kname, type) {
	instance = (
		"--image-id"
		$imageid
		"--key-name"
		$kname
		"--instance-type"
		$type
	)

	return $instance
}

fn aws_instance_setsecgrps(instance, secgrps) {
	instance <= append($instance, "--security-group-ids")
	instance <= append($instance, $secgrps)

	return $instance
}

fn aws_instance_setcount(instance, count) {
	instance <= append($instance, "--count")
	instance <= append($instance, $count)

	return $instance
}

fn aws_instance_setprivip(instance, privip) {
	instance <= append($instance, "--private-ip-address")
	instance <= append($instance, $privip)

	return $instance
}

fn aws_instance_setsubnet(instance, subnetid) {
	instance <= append($instance, "--subnet-id")
	instance <= append($instance, $subnetid)

	return $instance
}

fn aws_instance_setpubip(instance) {
	instance <= append($instance, "--associate-public-ip-address")

	return $instance
}

fn aws_instance_setprofile(instance, profile) {
	instance <= append($instance, "--iam-instance-profile")
	instance <= append($instance, $profile)

	return $instance
}

fn aws_instance_run(instance, tags) {
	instid <= (
		aws ec2 run-instances $instance |
		jq -r ".Instances[].InstanceId" |
		xargs echo -n
	)

	aws_tag($instid, $tags)

	return $instid
}

fn aws_instance_create(imgid, kname, secgrps, type, privip, subnetid, tags) {
	instid <= (
		aws ec2 run-instances
				--image-id $imgid
				--count 1
				--key-name $kname
				--security-group-ids $secgrps
				--instance-type $type
				--private-ip-address $privip
				--subnet-id $subnetid
				--associate-public-ip-address
				 |
		jq -r ".Instances[].InstanceId" |
		xargs echo -n
	)

	aws_tag($instid, $tags)

	return $instid
}

fn aws_instance_modify(resource, attr, value) {
	aws ec2 modify-instance-attribute --instance-id $resource --attribute $attr --value $value >[1=]
}

fn aws_instance_describe(filters) {
	filterStr = ""

	for f in $filters {
		if $filterStr == "" {
			filterStr = "Name="+$f[0]+",Values="+$f[1]
		} else {
			filterStr = $filterStr+",Name="+$f[0]+",Values="+$f[1]
		}
	}

	IFS = ("\n")

	instances <= (
		aws ec2 describe-instances
					--filters $filterStr |
		jq -j ".Reservations[].Instances[] | .InstanceId, \"  \", .Placement.AvailabilityZone, \"  \", .PrivateIpAddress, \"  \", .PublicIpAddress, \"\n\""
	)

	return $instances
}
