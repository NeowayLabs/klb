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
# The resource group can't exist and will be created during the process.
# This is important to ensure atomicity, to avoid any incomplete state to
# remain if something goes wrong with one of the copies the entire
# resource group is deleted with everything that is inside it returning
# to the original state. Resuming, or all the snapshots are copied or
# no snapshot is copied (there are no incomplete results).
#
# This function is intended to be used to copy snapshots to a
# different location from the original snapshots location and to do so
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
# An auxiliary temporary storage account will be created to do the work,
# we do our best to avoid leaking resources in the case of
# failures but azure is too awesome for us to guarantee that nothing
# will ever be leaked, so in case of errors it may be necessary
# to check if no resource has been leaked.
#
# The temporary storage account starts with the name "klbtmpsn".
fn azure_snapshot_copy(resgroup, location, snapshots_ids) {

    if azure_group_exists($resgroup) == "0" {
        return (), format("resgroup[%s] already exists", $resgroup)
    }

    fn log(msg, args...) {
        if len($args) == "0" {
            m = $msg
        } else {
            m <= format($msg, $args)
        }
        print("azure_snapshot_copy:%s\n", $m)
    }

    log("creating resgroup[%s]", $resgroup)
    azure_group_create($resgroup, $location)
    log("created resgroup with success")

    tmp_storage_acc <= _azure_snapshot_add_suffix("klbtmpsn")
    tmp_container = "klb-tmp-snapshot-copy"
    tmp_sku = "Premium_LRS"

    log("creating temporary storage account")
    err, storage_cleanup <= _azure_snapshot_create_storage_acc($tmp_storage_acc, $resgroup, $location, $tmp_sku, $tmp_container)

    fn cleanup() {
        err <= $storage_cleanup()
        if $err != "" {
            log("WARNING: error[%s] cleaning up storage account, resources may be leaked !!!", $err)
        }
    }

    fn err_cleanup() {
        cleanup()
        azure_group_delete($resgroup)
    }

    if $err != "" {
        err_cleanup()
        return format("error creating temporary storage account: %s", $err)
    }

    log("copied snapshots with success, cleaning up temporary storage account")
    cleanup()
    log("done")

    return (), "not implemented"
}

fn _azure_snapshot_create_storage_acc(acc, resgroup, location, sku, container) {

    fn nop() { return "" }

    fn log(msg, args...) {
        if len($args) == "0" {
            m = $msg
        } else {
            m <= format($msg, $args)
        }
        print("_azure_snapshot_create_storage_acc:%s\n", $m)
    }

    log("creating storage account: [%s] location: [%s]", $acc, $location)
    err <= azure_storage_account_create_storagev2($acc, $resgroup, $location, $sku)
    if $err != "" {
        return format("error creating storage account: %s",$err), $nop
    }
    log("created storage account with success")

    fn cleanup() {
        return azure_storage_account_delete($acc, $resgroup)       
    }

    log("getting storage account key")
    keys, err <= azure_storage_account_get_keys($acc, $resgroup)
    if $err != "" {
        return format("error getting storage account keys: %s", $err), $cleanup
    }

    if len($keys) == "0" {
        return format("error: unexpected empty keys when getting keys for account[%s]", $acc), $cleanup 
    }
    acc_triple = $keys[0]

    if len($acc_triple) != "3" {
        return format("error: returned account key is not a triple, instead it is: [%s]", $acc_triple), $cleanup
    }   

    acckey = $acc_triple[1]

    log("got storage account key: [%s]", $acckey)
    log("creating container: [%s]", $container)

    err <= azure_storage_container_create($container, $acc, $acckey)

    if $err != "" {
        return format("error creating container on storage account[%s]: %s", $acc, $err), $cleanup
    }

    log("successfully created storage account")
    return "", $cleanup
}

fn _azure_snapshot_add_suffix(name) {
	# Providing true uniqueness with the limits on the names is pretty hard :-)
	s <= head -n1 /dev/urandom | md5sum | tr -dc A-Za-z0-9 | cut -b 1-10
	return $name+$s
}