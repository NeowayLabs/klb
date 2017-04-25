# snapshot related functions

# azure_snapshot_create creates a new snapshot from the given
# srcid. The srcid can be a BLOB storage or even the disk of a
# running VM.
#
# Read more here: https://azure.microsoft.com/en-us/blog/azure-cli-managed-disks
#
# On success, returns the ID of the created snapshot.
# On failure it will explode in your face.
fn azure_snapshot_create(name, resgroup, srcid) {
	res <= (
		az snapshot create
				-g $resgroup
				-n $name
				--source $srcid |
		jq -r ".id"
	)

	return $res
}
