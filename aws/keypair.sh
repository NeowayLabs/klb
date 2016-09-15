# AWS Key Pair related functions

fn aws_keypair_create(name) {
	keyPair <= (
		aws ec2 create-key-pair
				--key-name kubernetes |
		jq -r ".KeyMaterial" |
		xargs echo -n a
	)

	return $keyPair
}

fn aws_keypair_delete(name) {
	-aws ec2 delete-key-pair --key-name $name

	return $status
}
