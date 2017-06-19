# Lock related functions

# azure_lock_create creates a new lock.
# Available lock types are: "CanNotDelete", "NotSpecified", "ReadOnly".
fn azure_lock_create(lockname, locktype, resgroup) {
	az lock create --lock-type $locktype --name $lockname --resource-group $resgroup
}

# azure_lock_delete deletes a lock with the given name.
fn azure_lock_delete(lockname, resgroup) {
	az lock delete --name $lockname --resource-type "locks" --namespace "Microsoft.Authorization" --resource-name $lockname --resource-group $resgroup
	# FIXME: delete is not sync, sometimes things may feel
	# implement a polling here perhaps ? Man..azure sucks :-)
}
