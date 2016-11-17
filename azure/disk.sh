# Disk related functions

fn azure_disk_attach_new(group, vm, storageaccount, size, name) {
	(
		azure vm disk attach-new --resource-group $group --vm-name $vm --storage-account-name $storageaccount --size-in-gb $size --vhd-name $name
	)
}

fn azure_disk_attach(group, vm, vhdurl) {
	(
		azure vm disk attach --resource-group $group --vm-name $vm --vhd-url $vhdurl
	)
}

fn azure_disk_detach(group, vm, lun) {
	(
		azure vm disk attach-new --resource-group $group --vm-name $vm --lun $lun
	)
}
