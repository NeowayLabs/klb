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

# azure_snapshot_list lists all available snapshots on
# the given resource group
#
# On success, returns a list of pairs (id, name),
# where name is the name of the snapshot and id is
# the id of the resource.
#
# On failure it will explode in your face (as usual).
fn azure_snapshot_list(resgroup) {
	# THIS IS COPIED FROM azure_vm_get_datadisks_ids_lun
	# It seems like a practical case for a zip function (like python)
	ids_names = ()

	snapshots     <= az snapshot list --resource-group $resgroup
	names_raw     <= echo $snapshots | jq -r ".[].name"
	ids_raw       <= echo $snapshots | jq -r ".[].id"
	ids           <= split($ids_raw, "\n")
	names         <= split($names_raw, "\n")
	size          <= len($ids)
	rangeend, err <= expr $size - 1

	if $err != "0" {
		return $ids_names
	}

	sequence <= seq "0" $rangeend
	range    <= split($sequence, "\n")

	for i in $range {
		idname = ($ids[$i] $names[$i])

		ids_names <= append($ids_names, $idname)
	}

	return $ids_names
}
