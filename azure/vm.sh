# Machine related functions

import klb/azure/group
import klb/azure/lock
import klb/azure/snapshot
import klb/azure/disk

# azure_vm_new creates a new instance of "virtual machine".
# `name` is the name of the virtual machine.
# `group` is name of resource group.
# `location` is the Azure Region.
fn azure_vm_new(name, group, location) {
	instance = (
		"--name"
		$name
		"--resource-group"
		$group
		"--location"
		$location
	)

	return $instance
}

# azure_vm_set_osdisk_id sets os disk ID of the virtual machine.
# If the this is called you should not set other OS parameters.
fn azure_vm_set_osdisk_id(instance, id) {
	instance <= append($instance, "--attach-os-disk")
	instance <= append($instance, $id)

	return $instance
}

# azure_vm_set_ostype sets ostype of "virtual machine".
# `instance` is the name of the instance.
# `ostype` is the type of OS installed on a custom VHD.
fn azure_vm_set_ostype(instance, ostype) {
	instance <= append($instance, "--os-type")
	instance <= append($instance, $ostype)

	return $instance
}

# azure_vm_set_vmsize sets the size of "Virtual Machine".
# `instance` is the name of the instance.
# `size` is the size of virtual machine.
# See https://goo.gl/vxW7we for size info.
fn azure_vm_set_vmsize(instance, vmsize) {
	instance <= append($instance, "--size")
	instance <= append($instance, $vmsize)

	return $instance
}

# azure_vm_set_customdata sets the custom init script file of "Virtual Machine".
# `instance` is the name of the instance.
# `customdata` is the location of custom init script file or text (cloud-config).
fn azure_vm_set_customdata(instance, customdata) {
	instance <= append($instance, "--custom-data")
	instance <= append($instance, $customdata)

	return $instance
}

# azure_vm_set_username sets the Username for the "Virtual Machine".
# `instance` is the name of the instance.
# `username` is the username for the virtual machine.
fn azure_vm_set_username(instance, username) {
	instance <= append($instance, "--admin-username")
	instance <= append($instance, $username)

	return $instance
}

# azure_vm_set_publickeyfile sets the SSH public key of "Virtual Machine".
# `instance` is the name of the instance.
# `publickeyfile` is the SSH public key or public key file path.
fn azure_vm_set_publickeyfile(instance, publickeyfile) {
	instance <= append($instance, "--ssh-key-value")
	instance <= append($instance, $publickeyfile)

	return $instance
}

# azure_vm_set_availset sets the availability set of "Virtual Machine".
# `instance` is the name of the instance.
# `availset` is the Name or ID of an existing availability set to add the
# Virtual Machine to.
fn azure_vm_set_availset(instance, availset) {
	instance <= append($instance, "--availability-set")
	instance <= append($instance, $availset)

	return $instance
}

# azure_vm_set_vnet sets the Virtual Network of "Virtual Machine".
# `instance` is the name of the instance.
# `vnet` is the name of the virtual network when creating a new one or
# referencing an existing one.
fn azure_vm_set_vnet(instance, vnet) {
	instance <= append($instance, "--vnet-name")
	instance <= append($instance, $vnet)

	return $instance
}

# azure_vm_set_subnet sets the Subnet of "Virtual Machine".
# `instance` is the name of the instance.
# `subnet` is the name of the subnet when creating a new VNet or
# referencing an existing one. Can also reference an existing
# subnet by ID.
fn azure_vm_set_subnet(instance, subnet) {
	instance <= append($instance, "--subnet")
	instance <= append($instance, $subnet)

	return $instance
}

# azure_vm_set_nic sets existing NIC to attach to the "Virtual Machine".
# `instance` is the name of the instance.
# `nicnames` is the name or ID of existing NIC to attach to the VM. The first
# NIC will be designated as primary. If omitted, a new 'NIC will be created. If
# an existing NIC is specified, do not specify subnet, vnet, public IP or NSG.
#
# Deprecated, use: azure_vm_set_nics
fn azure_vm_set_nic(instance, nic) {
	instance <= append($instance, "--nics")
	instance <= append($instance, $nic)

	echo "azure_vm_set_nic() is deprecated, use azure_vm_set_nics()"

	return $instance
}

# azure_vm_set_nics sets existing NICs to attach to the "Virtual Machine".
# `instance` is the name of the instance.
# `nicnames` is the names or IDs of existing NICs to attach to the VM. The first
# NIC will be designated as primary. If omitted, a new 'NIC will be created. If
# an existing NIC is specified, do not specify subnet, vnet, public IP or NSG.
fn azure_vm_set_nics(instance, nicnames) {
	fn join(list, sep) {
		out = ""

		for l in $list {
			out = $out+$l+$sep
		}

		out <= echo $out | sed "s/"+$sep+"$//g"

		return $out
	}

	nics     <= join($nicnames, " ")
	instance <= append($instance, "--nics")
	instance <= append($instance, $nics)

	return $instance
}

# azure_vm_set_storageaccount sets the storage account to "Virtual Machine".
# Note: Only applicable when use with 'azure_vm_set_unmanageddisk'
# `instance` is the name of the instance.
# `storageaccount` is the name to use when creating a new storage account or
# referencing an existing one.
fn azure_vm_set_storageaccount(instance, storageaccount) {
	instance <= append($instance, "--storage-account-name")
	instance <= append($instance, $storageaccount)

	return $instance
}

# azure_vm_set_osdiskvhd sets the OS disk name of "Virtual Machine".
# `instance` is the name of the instance.
# `osdiskvhd` is the name of the new Virtual Machine OS disk.
#
# Deprecated, use: azure_vm_set_osdiskname
fn azure_vm_set_osdiskvhd(instance, osdiskvhd) {
	instance <= append($instance, "--os-disk-name")
	instance <= append($instance, $osdiskvhd)

	echo "azure_vm_set_osdiskvhd() is deprecated, use azure_vm_set_osdiskname()"

	return $instance
}

# azure_vm_set_osdiskname sets the OS disk name of "Virtual Machine".
# `instance` is the name of the instance.
# `osdiskname` is the name of the new Virtual Machine OS disk.
fn azure_vm_set_osdiskname(instance, osdiskname) {
	instance <= append($instance, "--os-disk-name")
	instance <= append($instance, $osdiskname)

	return $instance
}

# azure_vm_set_datadisksize sets the empty managed data disk of "Virtual Machine".
# `instance` is the name of the instance.
# `datadisksize` is the space separated empty managed data disk sizes in GB.
fn azure_vm_set_datadisksize(instance, datadisksize) {
	instance <= append($instance, "--data-disk-sizes-gb")
	instance <= append($instance, $datadisksize)

	return $instance
}

# azure_vm_set_imageurn sets the image of "Virtual Machine".
# `instance` is the name of the instance.
# `imageurn` is the name of the operating system image (URN alias, URN,
# Custom Image name or ID, or VHD Blob URI)
fn azure_vm_set_imageurn(instance, imageurn) {
	instance <= append($instance, "--image")
	instance <= append($instance, $imageurn)

	return $instance
}

# azure_vm_set_storagesku sets the SKU storage account of "Virtual Machine".
# `instance` is the name of the instance.
# `storagesku` is the the sku of storage account to persist VM.
fn azure_vm_set_storagesku(instance, storagesku) {
	instance <= append($instance, "--storage-sku")
	instance <= append($instance, $storagesku)

	return $instance
}

# azure_vm_create creates a "Virtual Machine".
# `instance` is the name of the instance.
fn azure_vm_create(instance) {
	az vm create --output table $instance
}

fn azure_vm_delete(name, group) {
	(
		azure vm delete
			--name $name
			--resource-group $group
	)
}

fn azure_vm_get_ip_address(name, group, iface_index, ip_index) {
	info <= azure vm list-ip-address $group --json

	#echo $ips
	ip <= echo $info | jq -r ".[0].networkProfile.networkInterfaces["+$iface_index+"].expanded.ipConfigurations["+$ip_index+"].privateIPAddress"

	return $ip
}

# azure_vm_availset_create creates a new instance of Availset.
# `name` is the name of the Availset.
# `group` is name of resource group.
# `location` is the Azure Region.
fn azure_vm_availset_new(name, group, location) {
	instance = (
		"--name"
		$name
		"--resource-group"
		$group
		"--location"
		$location
	)

	return $instance
}

# azure_vm_availset_set_faultdomain sets Fault Domain of Availset.
# `instance` is the name of the instance.
# `count` is the Fault Domain count. Example: 2.
fn azure_vm_availset_set_faultdomain(instance, count) {
	instance <= append($instance, "--platform-fault-domain-count")
	instance <= append($instance, $count)

	return $instance
}

# azure_vm_availset_set_updatedomain sets Update Domain of Availset.
# `instance` is the name of the instance.
# `count` is the Update Domain count. Example: 2.
fn azure_vm_availset_set_updatedomain(instance, count) {
	instance <= append($instance, "--platform-update-domain-count")
	instance <= append($instance, $count)

	return $instance
}

# azure_vm_availset_set_unmanaged sets Contained VMs should use unmanaged disks.
# `instance` is the name of the instance.
fn azure_vm_availset_set_unmanaged(instance) {
	instance <= append($instance, "--unmanaged")

	return $instance
}

# azure_vm_availset_create creates a Availset.
# `instance` is the name of the instance.
fn azure_vm_availset_create(instance) {
	az vm availability-set create --output table $instance
}

# azure_vm_availset_delete deletes a Availset.
# `name` is the name of the Availset.
# `group` is name of resource group.
fn azure_vm_availset_delete(name, group) {
	az vm availability-set delete --output table --name $name --resource-group $group
}

# azure_vm_disk_attach attaches an existing disk to the VM.
fn azure_vm_disk_attach(name, resgroup, diskID) {
	az vm disk attach -g $resgroup --vm-name $name --disk $diskID
}

# azure_vm_disk_attach_lun does the same as azure_vm_disk_attach
# but with the specificied LUN.
fn azure_vm_disk_attach_lun(name, resgroup, diskID, lun) {
	az vm disk attach -g $resgroup --vm-name $name --disk $diskID --lun $lun
}

# azure_vm_disk_attach_new creats a new disk and attaches to the VM.
fn azure_vm_disk_attach_new(name, resgroup, diskname, size, sku) {
	az vm disk attach -g $resgroup --vm-name $name --disk $diskname --new --size-gb $size --sku $sku
}

# azure_vm_get_datadisks_ids will returns a list with the
# id's of all the managed data disks of the given VM.
# If the VM has no data disk it will return an empty list.
# The id of the OSDisk will not be returned.
#
# These ID's are well suited to be used on snapshot creation.
fn azure_vm_get_datadisks_ids(name, resgroup) {
	info    <= azure_vm_get_rawinfo($name, $resgroup)
	ids_raw <= echo $info | jq -r ".storageProfile.dataDisks[].managedDisk.id"
	ids     <= split($ids_raw, "\n")

	return $ids
}

# azure_vm_get_datadisks_ids_lun will returns a list
# of pairs (diskid, disklun).
#
# The LUN is only useful if you are trying to replicate
# a VM exact state, so you need to know the disks LUN so
# they will be attached on the new VM with the same device names.
#
# These ID's are the same returned by azure_vm_get_datadisks_ids
fn azure_vm_get_datadisks_ids_lun(name, resgroup) {
	ids_luns  = ()
	info      <= azure_vm_get_rawinfo($name, $resgroup)
	disks_raw <= echo $info | jq -r ".storageProfile.dataDisks[]"
	ids_raw <= echo $disks_raw | jq -r ".managedDisk.id"

	if $ids_raw == "" {
		echo "no datadisks found for vm: " + $name
		return $ids_luns
	}
	ids       <= split($ids_raw, "\n")

	luns_raw  <= echo $disks_raw | jq -r ".lun"
	if $luns_raw == "" {
		echo "no datadisks found for vm: " + $name
		return $ids_luns
	}

	luns      <= split($luns_raw, "\n")
	size      <= len($ids)

	rangeend  <= expr $size - 1
	sequence  <= seq "0" $rangeend
	range     <= split($sequence, "\n")

	for i in $range {
		idlun = ($ids[$i] $luns[$i])

		ids_luns <= append($ids_luns, $idlun)
	}

	return $ids_luns
}

# azure_vm_get_osdisk_id will return the osdisk ID.
#
# This ID is well suited to be used on snapshot creation.
fn azure_vm_get_osdisk_id(name, resgroup) {
	info <= azure_vm_get_rawinfo($name, $resgroup)
	id   <= echo $info | jq -r ".storageProfile.osDisk.managedDisk.id"
	if $id == "" {
		return "", format(
			"unable to find managed OS disk on vm[%s] resgroup[%s], vm probably do not exist",
			$name,
			$resgroup,
		)
	}
	return $id, ""
}

# azure_vm_get_rawinfo will return the raw encoded JSON data with
# all the VM info.
fn azure_vm_get_rawinfo(name, resgroup) {
	info <= az vm show --resource-group $resgroup --name $name

	return $info
}

# azure_vm_get_disks_ids will return the id of all disks on the VM.
#
# It will be a list including the osdisk id and the datadisks id's.
fn azure_vm_get_disks_ids(name, resgroup) {
	osdiskid, err <= azure_vm_get_osdisk_id($name, $resgroup)
	if $err != "" {
		return (), $err
	}
	datadisks <= azure_vm_get_datadisks_ids($name, $resgroup)
	disks     <= append($datadisks, $osdiskid)

	return $disks, ""
}

# azure_vm_stop stops a running VM
fn azure_vm_stop(vmname, resgroup) {
	az vm stop -g $resgroup -n $vmname
}

# azure_vm_start starts a stopped VM
fn azure_vm_start(vmname, resgroup) {
	az vm start -g $resgroup -n $vmname
}

# azure_vm_backup_create will create a full backup from the given VM.
# The backup is created following a naming convention for the resources
# that it creates, this enables the azure_vm_backup_recover function to work.
# Conceptually we are encoding information required for proper recover (metadata)
# on the name of the resources, so we don't need a third party storage.
#
# The "prefix" parameter gives you a namespace that you can use to organize
# backups of different applications on the same subscription. This namespace
# is built by appending the provided string as a prefix on the name of
# the resource group that will be created to hold the backup.
#
# When you call this function, the first step is to create a resource
# group with the name following this pattern:
#
# <prefix>-klb-backup-<timestamp>-<vmname>
#
# Where the prefix is the one you passed as a parameter.
# If there is already a resource group with this name, the
# creation will fail. It is paramount to the proper work of
# the backup functions that the ONLY thing inside the resource
# group are the VM disks snapshots.
#
# The <timestamp> will follow this pattern:
#
# <year>.<month>.<day>.<hour>
#
# Calling: azure_vm_backup_create("test", "testgroup", "staging")
#
# Could (timestamp may vary) create the resource group:
#
# "staging-bkp-2017.05.28.1930-test"
#
# The resource group name convention is important to be know since
# you must manage these resource groups and delete them.
# There will be also a convention on how snapshots are stored inside
# this resource group but they are not documented and you should not
# rely on them, they are implementation details.
#
# Backup resource groups should never be changed, because of that
# they are read only locked after all snapshots are added. There is
# also a delete lock to avoid deleting backups on accident.
# The azure_vm_backup_delete function will release the locks and delete
# a backup resource group for you.
#
# During the backup procedure the VM must be stopped (you can use the
# azure_vm_stop function to do that). It is a programming error to call
# this function with a VM that is running as a parameter, since it is
# not safe to take snapshots from a running VM.
#
# After the function returns it is safe to start the VM, all the snapshots
# have been taken.
#
# Be aware that resource group names have a lame limit of 64
# characters. Since we need to create locks the final limit will
# be 60 characters total (including the timestamp).
# So avoid long names for VM's and prefixes.
#
# On success it will return the name of the created resource group and
# an empty string as error. On error it will return an empty string as resource
# group and a non-empty error string with details on the failure.
fn azure_vm_backup_create(vmname, resgroup, prefix, location) {
	timestamp <= date "+%Y.%m.%d.%H%M"
	# WHY: We need some chars for the lock names,
	# based on the resgroup name.
	max_resgroup_size = "60"

	bkp_resgroup = $prefix+"-bkp-"+$timestamp+"-"+$vmname

	bkp_resgroup_len <= len($bkp_resgroup)
	_, err <= test $max_resgroup_size -gt $bkp_resgroup_len
	if $err != "0" {
		return "", format("error: resgroup name %q is too bigger than %q", $bkp_resgroup, $max_resgroup_size)
	}

	if azure_group_exists($bkp_resgroup) == "0" {
		return "", format("error: resource group already exists: %q", $bkp_resgroup)
	}

	echo "vm.backup.create: getting VM disks IDs"
	echo "vm.backup.create: vm name: " + $vmname
	echo "vm.backup.create: resgroup: " + $resgroup

	osdiskid, err <= azure_vm_get_osdisk_id($vmname, $resgroup)
	if $err != "" {
		return "", $err
	}
	echo "got os disk id: " + $osdiskid

	disks_ids_luns <= azure_vm_get_datadisks_ids_lun($vmname, $resgroup)
	echo "vm.backup.create: creating resource group: "+$bkp_resgroup
	echo "vm.backup.create: at location: "+$location

	azure_group_create($bkp_resgroup, $location)


	# WHY: name used later on the recover phase, do NOT change this
	# unless you are absolutely SURE of what you are doing
	snapshot_name <= _azure_vm_backup_get_osdisk_name()

	echo "vm.backup.create: creating OS disk snapshot: "+$snapshot_name+" from disk id: "+$osdiskid

	snapshotid <= azure_snapshot_create($snapshot_name, $bkp_resgroup, $osdiskid)

	echo "vm.backup.create: created snapshot id: "+$snapshotid

	for idlun in $disks_ids_luns {
		id  = $idlun[0]
		lun = $idlun[1]

		# WHY: Encode lun on the name as metadata, use it later to restore
		# Change this and the whole world will collapse :-)
		snapshot_name <= _azure_vm_backup_datadisk_name($lun)

		echo "vm.backup.create: creating datadisk snapshot: "+$snapshot_name+" from disk id: "+$id

		snapshotid <= azure_snapshot_create($snapshot_name, $bkp_resgroup, $id)

		echo "vm.backup.create: created snapshot id: "+$snapshotid
	}

	echo "vm.backup.create: backup finished with success, creating lock"

	nodelete <= _azure_vm_backup_get_nodelete_lock($bkp_resgroup)

	azure_lock_create($nodelete, "CanNotDelete", $bkp_resgroup)

	readonly <= _azure_vm_backup_get_readonly_lock($bkp_resgroup)

	azure_lock_create($readonly, "ReadOnly", $bkp_resgroup)

	echo "vm.backup.create: created lock, finished with success"

	return $bkp_resgroup, ""
}

# azure_vm_backup_list returns the list of all backups
# for the given vm name + prefix. They are the same parameters
# you used to create the backup.
#
# The backup list is a list of resource groups names, where each
# resource group is a backup.
#
# The list will be ordered, from the more recent to the oldest backup.
fn azure_vm_backup_list(vmname, prefix) {
	resgroups <= azure_group_get_names()

	filtered = ""

	for resgroup in $resgroups {
		hasprefix, _ <= echo $resgroup | grep "^"+$prefix
		hasvmname, status <= echo $hasprefix | grep $vmname+"$"

		if $status == "0" {
			filtered = $filtered+$resgroup+"\n"
		}
	}

	return _azure_vm_backup_order_list($filtered)
}

# azure_vm_backup_list_all returns the list of all backups
# for the given prefix. If there is backups for multiple VM's
# for the given prefix it will return all of them.
#
# The return value is the same as azure_vm_backup_list, just aggregating
# results for all VMs instead of a single one.
fn azure_vm_backup_list_all(prefix) {
	resgroups <= azure_group_get_names()

	filtered = ""

	for resgroup in $resgroups {
		hasprefix, status <= echo $resgroup | grep "^"+$prefix

		if $status == "0" {
			filtered = $filtered+$resgroup+"\n"
		}
	}

	return _azure_vm_backup_order_list($filtered)
}

# azure_vm_backup_delete deletes a backup. This function
# will also remove the locks that prevents backups deletion.
fn azure_vm_backup_delete(backup_resgroup) {

	dellock <= _azure_vm_backup_get_nodelete_lock($backup_resgroup)
	readlock <= _azure_vm_backup_get_readonly_lock($backup_resgroup)

	echo "backup delete: resource group: " + $backup_resgroup
	echo "backup delete: removing lock: " + $dellock
	azure_lock_delete($dellock, $backup_resgroup)
	echo "backup delete: removing lock: " + $readlock
	azure_lock_delete($readlock, $backup_resgroup)

	# WHY: Deleting locks is not synchronous, this is a VERY
	# lame workaround while we don't fix the azure_lock_delete function.
	lame_workaround_sleep = "5"
	sleep $lame_workaround_sleep
	echo "backup delete: locks removed, deleting resource group: " + $backup_resgroup
	azure_group_delete($backup_resgroup)
}

# azure_vm_backup_recover will recover a previously generated
# backup. The backup resource group must have been generated using
# azure_vm_backup_create, since a whole convention on resource
# naming will be required in the recovery process.
#
# The vminstance parameter is a instance of the vm object,
# created with azure_vm_new and configured just like you
# would do to create a new VM.
#
# The main differences is
# that no os disk should be configured, since the os disk and
# datadisks will be obtained from the backup_resgroup, and
# no storage-sku should be set on the vm instance, since it
# will be defined on the disks created from the backup.
#
# The backup_resgroup is the name of the resource group
# where the snapshots are stored just as it is returned by
# azure_vm_backup_create.
#
# It is an error to provide a vm instance that has not
# a os type setted using azure_vm_set_ostype.
#
# This function returns an empty string on success or a
# non empty error message if it fails.
fn azure_vm_backup_recover(instance, storagesku, backup_resgroup) {

	fn log(msg) {
		echo "vm.backup.recover: " + $msg
	}

	log("getting info from vm")
	resgroup <= _azure_vm_get($instance, "resource-group")
	location <= _azure_vm_get($instance, "location")
	vmname <= _azure_vm_get($instance, "name")
	ostype <= _azure_vm_get($instance, "os-type")

	log("vm name: " + $vmname)
	log("vm resgroup: " + $resgroup)
	log("vm location: " + $location)
	log("vm os type: " + $ostype)

	if $vmname == "" {
		return "unable to get the 'name' from the given vm instance"
	}

	if $resgroup == "" {
		return "unable to get the 'resource-group' from the given vm instance"
	}

	if $location == "" {
		return "unable to get the 'location' from the given vm instance"
	}

	if $ostype == "" {
		return "unable to get the 'os-type' from the given vm instance"
	}

	osdiskname <= _azure_vm_get($instance, "os-disk-name")
	if $osdiskname != "" {
		msg <= format("found os disk name %q on vm instance", $osdiskname)
		return $msg + "should not call azure_vm_set_osdiskname on a vm that is being recovered from backup"
	}

	sku <= _azure_vm_get($instance, "storage-sku")
	if $sku != "" {
		msg <= format("found storage-sku %q on vm instance", $sku)
		return $msg + "should not call azure_vm_set_storagesku on a vm that is being recovered from backup"
	}

	log("loading snapshots from backup: " + $backup_resgroup)
	snapshots <= azure_snapshot_list($backup_resgroup)
	log(format("loaded snapshots, parsing results: %s", $snapshots))
	osdiskid = ""
	datadisks = ()

	osdiskname <= _azure_vm_backup_get_osdisk_name()
	for snapshot in $snapshots {
		log(format("parsing: %s", $snapshot))
		id = $snapshot[0]
		name = $snapshot[1]
		if $name == $osdiskname {
			osdiskid = $id
		} else {
			lun <= _azure_vm_backup_datadisk_lun($name)
			idlun = ($id $lun)
			datadisks <= append($datadisks, $idlun)
		}
	}

	log("os disk id:[" + $osdiskid + "]")
	if $osdiskid == "" {
		return format(
			"unable to find osdisk id on backup resource group: %q, corrupted backup ?",
			$backup_resgroup,
		)
	}

	log("creating os disk")
	osdiskname = $vmname + "-osdisk"
	d <= azure_disk_new($osdiskname, $resgroup, $location)
	d <= azure_disk_set_source($d, $osdiskid)
	d <= azure_disk_set_sku($d, $storagesku)
	osdisk <= azure_disk_create($d)

	log("created os disk: " + $osdisk)
	instance <= azure_vm_set_osdisk_id($instance, $osdisk)

	log("creating VM")
	azure_vm_create($instance)
	log("created VM, stopping it so we can attach disks")
	# https://feedback.azure.com/forums/216843-virtual-machines/suggestions/6750456-allow-to-create-vm-without-starting-it-immediatell
	azure_vm_stop($vmname, $resgroup)

	log("attaching datadisks")
	for datadisk in $datadisks {
		id = $datadisk[0]
		lun = $datadisk[1]

		log("creating disk from snapshot: " + $id)
		log("disk will have LUN: " + $lun)

		diskname = $vmname + "-disk-" + $lun
		d <= azure_disk_new($diskname, $resgroup, $location)
		d <= azure_disk_set_source($d, $id)
		d <= azure_disk_set_sku($d, $storagesku)
		diskid <= azure_disk_create($d)

		log("created disk id: " + $diskid)
		log("attaching on VM")
		azure_vm_disk_attach($vmname, $resgroup, $diskid)
		log("attached")
	}

	log("starting VM with all disks attached")
	azure_vm_start($vmname, $resgroup)
	log("finished recover with success")

	return ""
}

fn _azure_vm_backup_get_nodelete_lock(bkp_resgroup) {
	return "del-"+$bkp_resgroup
}

fn _azure_vm_backup_get_readonly_lock(bkp_resgroup) {
	return "ro-"+$bkp_resgroup
}

fn _azure_vm_backup_order_list(backup_resgroups) {
	ordered_raw <= echo $backup_resgroups | sort -r
	ordered <= split($ordered_raw, "\n")
	res = ()
	# WHY: handle possible trailing newlines
	for o in $ordered {
		if $o != "" {
			res <= append($res, $o)
		}
	}
	return $res
}

fn _azure_vm_backup_get_osdisk_name() {
	return "osdisk"
}

fn _azure_vm_backup_datadisk_name(lun) {
	return "datadisk-"+$lun
}

fn _azure_vm_backup_datadisk_lun(name) {
	tokens <= split($name, "-")
	if len($tokens) != "2" {
		echo "invalid backup datadisk name: " + $name
		exit("1")
	}
	return $tokens[1]
}

fn _azure_vm_get(instance, cfgname) {
	cfgname = "--" + $cfgname
	size      <= len($instance)
	rangeend  <= expr $size "-" "2"
	sequence  <= seq "0" $rangeend
	range     <= split($sequence, "\n")

	ids_names = ()

	for i in $range {
		cfgval_index <= expr $i "+" "1"
		name = $instance[$i]
		if $name == $cfgname {
			return $instance[$cfgval_index]
		}
	}

	return ""
}
