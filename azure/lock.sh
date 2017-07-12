# Lock related functions

# azure_lock_create creates a new lock.
# Available lock types are: "CanNotDelete", "NotSpecified", "ReadOnly".
fn azure_lock_create(lockname, locktype, resgroup) {
	az lock create --lock-type $locktype --name $lockname --resource-group $resgroup
}

# azure_lock_delete deletes a lock with the given name.
fn azure_lock_delete(lockname, resgroup) {
	err <= _azure_lock_delete($lockname, $resgroup, "0")
	return $err
}

fn _azure_lock_delete(lockname, resgroup, trycount) {
	
	maxretry = "10"

	fn error(reason) {
		return format(
			"azure_lock_delete: unable to delete lock[%s] resgroup[%s]: reason: %s",
			$lockname,
			$resgroup,
			$reason,
		)
	}

	_, status <= test $trycount "-lt" $maxretry
	if $status != "0" {
		return error("exceeded max retry: " + $maxretry)
	}

	fn log(msg) {
		print("azure_lock_delete: %s\n", $msg)
	}

	# WHY: delete is not sync, sometimes things may fail silently
	# implement a polling for now.
	log(format("deleting lock[%s] from resgroup[%s]", $lockname, $resgroup))

	out, status <= az lock delete --name $lockname --resource-type "locks" --namespace "Microsoft.Authorization" --resource-name $lockname --resource-group $resgroup
	if $status != "0" {
		return error(format("unable to delete lock, error details: [" + $out + "]"))
	}

	log(format("deleted lock[%s] successfully, checking if it still exists", $lockname))

	out, status <= az lock list | grep $lockname
	if $status == "0" {
		log(format("lock[%s] still exists, trying again", $lockname))
		pollingtimesec = "1"
		trycount <= expr $trycount "+" "1"
		sleep $pollingtimesec
		return _azure_lock_delete($lockname, $resgroup, $trycount)
	}

	log(format("lock[%s] not found, success", $lockname))
	return ""
}
