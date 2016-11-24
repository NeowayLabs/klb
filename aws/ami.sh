# AMI related functions

fn aws_ami_filter(filters) {
	filterStr = ""

	for f in $filters {
		if $filterStr == "" {
			filterStr = "Name="+$f[0]+",Values="+$f[1]
		} else {
			filterStr = $filterStr+",Name="+$f[0]+",Values="+$f[1]
		}
	}

	amis <= (
		aws ec2 describe-images
				--filters $filterStr |
		jq ".Images[]"
	)

	return $amis
}

fn aws_ami_get(imgid) {
	ami <= aws ec2 describe-images --image-ids $imgid | jq ".Images[0]"

	return $ami
}
