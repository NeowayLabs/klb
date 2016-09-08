# Elastic IP address

# Create an AWS Elastic IP address
fn aws_eip_create(tags) {
	eipId <= (
		aws ec2 allocate-address --domain vpc |
		jq ".AllocationId"
	)

        aws_tags($eipId, $tags)
        return $eipId
}

# Associate an EIP to an instance
fn aws_eip_associate(eipid, instanceId) {
	assocId <= (
		aws ec2 associate-address
				--instance-id $instanceId
				--public-ip $eipId |
		jq ".AssociationId"
	)

        return $assocId
}
