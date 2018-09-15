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
	out, status <= (
		az snapshot create
				-g $resgroup
				-n $name
				--sku $sku
				--source $srcid
	)

    if $status != "0" {
        return "", format("error creating snapshot: %s", $out)
    }

    res <= echo $out | jq -r ".id"
	return $res, ""
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
fn azure_snapshot_copy(resgroup, location, sku, snapshots_ids) {

    if azure_group_exists($resgroup) == "0" {
        return (), format("resgroup[%s] already exists", $resgroup)
    }

    fn log(msg, args...) {
        print("azure_snapshot_copy:%s\n", format($msg, $args...))
    }

    log("getting access rights to the snapshots")
    # Supposing it wont take more than a day to copy snapshots
    sas_timeout_sec = "86400"
    snapshots_ids_sas, err <= azure_snapshot_grant_access($snapshots_ids, $sas_timeout_sec)
    if $err != "" {
        return (), format("error granting access to snapshots: %s", $err)
    } 

    log("creating resgroup[%s]", $resgroup)
    azure_group_create($resgroup, $location)
    log("created resgroup with success")

    tmp_storage_acc <= _azure_snapshot_add_suffix("klbtmpsn")
    tmp_container = "klb-tmp-snapshot-copy"

    log("creating temporary storage account")
    tmp_acckey, storage_cleanup, err <= _azure_snapshot_create_storage_acc($tmp_storage_acc, $resgroup, $location, $sku, $tmp_container)

    fn cleanup() {
        err <= $storage_cleanup()
        if $err != "" {
            log("WARNING: error[%s] cleaning up storage account, resources may be leaked !!!", $err)
        } else {
            log("removed temporary storage account with success")
        }
    }

    fn err_cleanup() {
        log("error ocurred, deleting all temporary resources and resource group")
        cleanup()
        azure_group_delete($resgroup)
        log("error cleanup succeeded")
    }

    if $err != "" {
        err_cleanup()
        return (), format("error creating temporary storage account: %s", $err)
    }

    log("created temporary storage account with success, account key:[%s]", $tmp_acckey)

    tmp_blobs = ()

    for snapshot_id_sas in $snapshots_ids_sas {
        if len($snapshot_id_sas) != "2" {
             err_cleanup()
             return (), format("error: expected pair (snapshotid, snapshot_sas) but got [%s]", $snapshot_id_sas)
        }

        snapshot_id = $snapshot_id_sas[0]
        snapshot_sas = $snapshot_id_sas[1]
        log("starting copy of snapshot[%s] sas[%s]", $snapshot_id, $snapshot_sas)

        snapshot_name, err <= azure_snapshot_name($snapshot_id)
        if $err != "" {
            err_cleanup()
            return (), format("error getting snapshot[%s] name: %s", $snapshot_id, $err)
        }
        tmpblob <= _azure_snapshot_add_suffix($snapshot_name + "-")
        log("generated unique blob name[%s] for snapshot named[%s] with id[%s]", $tmpblob, $snapshot_name, $snapshot_id)

        copyid, err <= azure_storage_blob_copy_start($tmp_storage_acc, $tmp_acckey, $tmp_container, $tmpblob, $snapshot_sas)

        if $err != "" {
            err_cleanup()
            return (), format("error starting snapshot async copy to another location: %s", $err)
        }

        log("started async copy to another location with success, op id[%s]", $copyid)
        # FIXME: CANT FIND THE URL ANYWHERE =/
        tmpblob_url <= format("https://%s.blob.core.windows.net/%s/%s", $tmp_storage_acc, $tmp_container, $tmpblob)
        log("generated temporary blob url: [%s]", $tmpblob_url)

        tmp_blobs <= append($tmp_blobs, ($tmpblob $tmpblob_url $copyid $snapshot_name))
    }

    err <= _azure_snapshot_wait_blobs_copy($tmp_storage_acc, $tmp_acckey, $tmp_container, $tmp_blobs)
    if $err != "" {
        err_cleanup()
        return (), format("error waiting for blobs to finish copying: %s", $err)
    }

    log("copied all snapshots to desired location, starting to create snapshots")

    res = ()

    for tmp_blob in $tmp_blobs {
        blob_name = $tmp_blob[0]
        blob_url  = $tmp_blob[1]
        snapshot_name = $tmp_blob[3]
 
        log("creating new snapshot[%s] at location[%s] from blob[%s] url[%s]", $snapshot_name, $location, $blob_name, $blob_url)

        copied_snapshot_id, err <= azure_snapshot_create($snapshot_name, $resgroup, $blob_url, $sku)
        if $err != "" {
            err_cleanup()
            return "", format("error copying snapshot: %s", $err)
        }
        log("copied snapshot id: [%s]", $copied_snapshot_id)
        res <= append($res, $copied_snapshot_id)
    }

    log("copied snapshots with success, cleaning up temporary storage account")
    cleanup()
    log("done")

    return $res, ""
}

fn azure_snapshot_name(snapshot_id) {
    out, status <= az snapshot show --ids $snapshot_id
    if $status != "0" {
        return "", format("error getting snapshot name: %s", $out)
    }
    name <= echo $out | jq -r ".name"
    return $name, ""
}

# azure_snapshot_grant_access will get a list of snapshot ids and will return a
# list of pairs of (snapshot_id snapshot_sas) where the sas is the read access
# that has been granted to read the snapshot (useful for operations like a
# blob storage copy).
#
# In success it return the list and an empty error string, otherwise it
# returns an empty list and the error message.
fn azure_snapshot_grant_access(snapshots_ids, sas_timeout_sec) {
    
    fn log(msg, args...) {
        print("azure_snapshot_grant_access:%s\n", format($msg, $args...))
    }

    res = ()

    for snapshot_id in $snapshots_ids {
        sas, status <= az snapshot grant-access --ids $snapshot_id --duration-in-seconds $sas_timeout_sec --query "[accessSas]" -o tsv
        
        if $status != "0" {
            return (), format("error granting read access to snapshot[%s]: %s", $snapshot_id, $sas)
        }
    
        log("got read access granted with sas[%s] for snapshot id [%s]", $sas, $snapshot_id)
        res <= append($res, ($snapshot_id $sas))
    }

    return $res, ""
}

fn _azure_snapshot_wait_blobs_copy(account, acckey, container, blobs) {

    # TODO: could avoid this loggers duplication with a logger creation function
    fn log(msg, args...) {
        print("_azure_snapshot_wait_blobs_copy:%s\n", format($msg, $args...))
    }

    # WHY: because all good algorithms starts with a while true :D
    # FIXME: seriously though, it would be good to have a timeout here x_x
    for {
        completed = ()

        for blob in $blobs {
            blob_name = $blob[0]
            copyid = $blob[2]

            output, status <= az storage blob show --container-name $container --name $blob_name --account-key $acckey --account-name $account

            if $status != "0" {
                return format("error getting status of blob[%s] copy operation", $blob_name)
            }

            copystatus <= echo $output | jq -r ".properties.copy.status"
            copyprogress <= echo $output | jq -r ".properties.copy.progress"
            got_copyid <= echo $output | jq -r ".properties.copy.id"

            if $got_copyid != $copyid {
                return format("expected copyid[%s] but got[%s] copying blob[%s]", $copyid, $got_copyid, $blob_name)
            }

            log("blob[%s] progress[%s] status[%s]", $blob_name, $copyprogress, $copystatus)
            if $copystatus == "success" {
                completed <= append($completed, $blob_name)
            }
        }

        if len($completed) == len($blobs) {
            return ""
        }

        log("there are still copies going on, blobs completed: [%s]", $completed)
    }

}

fn _azure_snapshot_create_storage_acc(acc, resgroup, location, sku, container) {

    fn nop() { return "" }

    fn log(msg, args...) {
        print("_azure_snapshot_create_storage_acc:%s\n", format($msg, $args...))
    }

    log("creating storage account: [%s] location: [%s]", $acc, $location)
    err <= azure_storage_account_create_storagev2($acc, $resgroup, $location, $sku)
    if $err != "" {
        return "", $nop, format("error creating storage account: %s",$err)
    }
    log("created storage account with success")

    fn cleanup() {
        return azure_storage_account_delete($acc, $resgroup)       
    }

    log("getting storage account key")
    keys, err <= azure_storage_account_get_keys($acc, $resgroup)
    if $err != "" {
        return "", $cleanup, format("error getting storage account keys: %s", $err)
    }

    if len($keys) == "0" {
        return "", $cleanup, format("error: unexpected empty keys when getting keys for account[%s]", $acc)
    }
    acc_triple = $keys[0]

    if len($acc_triple) != "3" {
        return "", $cleanup, format("error: returned account key is not a triple, instead it is: [%s]", $acc_triple)
    }   

    acckey = $acc_triple[1]

    log("got storage account key: [%s]", $acckey)
    log("creating container: [%s]", $container)

    err <= azure_storage_container_create($container, $acc, $acckey)

    if $err != "" {
        return "", $cleanup, format("error creating container on storage account[%s]: %s", $acc, $err)
    }

    log("successfully created storage account")
    return $acckey, $cleanup, ""
}

fn _azure_snapshot_add_suffix(name) {
	# Providing true uniqueness with the limits on the names is pretty hard :-)
	s <= head -n1 /dev/urandom | md5sum | tr -dc A-Za-z0-9 | cut -b 1-10
	return $name+$s
}
