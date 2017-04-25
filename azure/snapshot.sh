# snapshot related functions

# azure_snapshot_create creates a new snapshot from the given
# srcid. The srcid can be a BLOB storage or even the disk of a
# running VM.
#
# Read more here: https://azure.microsoft.com/en-us/blog/azure-cli-managed-disks
fn azure_snapshot_create(name, resgroup, srcid) {
	az snapshot create -g $resgroup -n $name --source $srcid
}
