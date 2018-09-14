# snapshot related functions

# azure_snapshot_create creates a new snapshot from the given
# srcid. The srcid can be a BLOB storage or even the disk of a
# running VM. The provided sku will be used to define the type
# of the snapshot disk.
#
# Read more here: https://azure.microsoft.com/en-us/blog/azure-cli-managed-disks
#
# On success, returns the ID of the created snapshot.
# On failure it will explode in your face.
fn azure_snapshot_create(name, resgroup, srcid, sku) {
	res <= (
		az snapshot create
				-g $resgroup
				-n $name
				--sku $sku
				--source $srcid |
		jq -r ".id"
	)

	return $res
}

# azure_snapshot_list lists all available snapshots on
# the given resource group
#
# On success, returns a list of pairs (id, name),
# where name is the name of the snapshot and id is
# the id of the resource and an empty error string.
#
# On failure it will return an empty list and an error
# string with details on the failure.
fn azure_snapshot_list(resgroup) {
	# THIS IS COPIED FROM azure_vm_get_datadisks_ids_lun
	# It seems like a practical case for a zip function (like python)
	ids_names = ()

	snapshots, status <= az snapshot list --resource-group $resgroup

	if $status != "0" {
		return (), format("unable to list snapshots on resource group[%s]", $resgroup)
	}

	names_raw, status <= echo $snapshots | jq -r ".[].name"

	if $status != "0" {
		return (), format("error parsing snapshot names from[%s]", $snapshots)
	}

	ids_raw, status <= echo $snapshots | jq -r ".[].id"

	if $status != "0" {
		return (), format("error parsing snapshot ids from[%s]", $snapshots)
	}

	ids         <= split($ids_raw, "\n")
	names       <= split($names_raw, "\n")
	size        <= len($ids)
	rangeend, _ <= expr $size - 1
	sequence    <= seq "0" $rangeend
	range       <= split($sequence, "\n")

	for i in $range {
		idname = ($ids[$i] $names[$i])

		ids_names <= append($ids_names, $idname)
	}

	return $ids_names, ""
}

# azure_snapshot_copy will copy the given list of snapshots ids
# to the resource group and location informed.
#
# This function is intended to be used to copy snapshots to a
# different location from the original snapshots and to do so
# it will do a lot of work that can take some time.
#
# At the time this is being written there is not support to
# copying or generating snapshots from a disk to a location different
# from the original disk/snapshot:
#
# - https://stackoverflow.com/questions/47759200/creating-a-managed-disk-from-snapshot-in-different-region-azure
# - https://docs.microsoft.com/en-us/azure/virtual-machines/scripts/virtual-machines-linux-cli-sample-copy-snapshot-to-storage-account
#
# This function will do this job for you. On success it returns a list of
# the created snapshots and an empty error. On failure it will return an empty
# list and a string error with details.
#
# Auxiliary temporary storage accounts will be created to do the work,
# we do our best to avoid leaking these resources in the case of
# failures but azure is too awesome for us to guarantee that nothing
# will ever be leaked, so in case of errors it may be necessary
# to check if no resource has been leaked.
#
# All temporary resources starts with the name "klb-tmp-sn".
fn azure_snapshot_copy(resgroup, location, snapshots_ids) {
    return (), "not implemented"
}
