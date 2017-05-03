# Lock related functions

# azure_lock_create creates a new lock.
# Available lock types are: "CanNotDelete", "NotSpecified", "ReadOnly".
fn azure_lock_create(lockname, locktype, resgroup) {
	az lock create --lock-type $locktype --name $lockname --resource-group $resgroup
}

# azure_lock_delete deletes a lock with the given name.
fn azure_lock_delete(lockname) {
	az lock delete --name $lockname
}
