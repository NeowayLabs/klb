# Machine related functions

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

	nics     <= join($nicnames, ",")
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
	ids_raw <= (
		az vm show
			--resource-group $resgroup
			--name $name |
		jq -r ".storageProfile.dataDisks[].managedDisk.id"
	)

	ids     <= split($ids_raw, "\n")

	return $ids
}

# azure_vm_get_osdisk_id will return the osdisk ID.
#
# This ID is well suited to be used on snapshot creation.
fn azure_vm_get_osdisk_id(name, resgroup) {
	id <= (
		az vm show
			--resource-group $resgroup
			--name $name |
		jq -r ".storageProfile.osDisk.managedDisk.id"
	)

	return $id
}

# azure_vm_get_disks_ids will return the id of all disks on the VM.
#
# It will be a list including the osdisk id and the datadisks id's.
fn azure_vm_get_disks_ids(name, resgroup) {
	osdiskid  <= azure_vm_get_osdisk_id($name, $resgroup)
	datadisks <= azure_vm_get_datadisks_ids($name, $resgroup)
	disks     <= append($datadisks, $osdiskid)

	return $disks
}

# azure_vm_backup_create will create a full backup from the given VM.
# The backup is created following a naming pattern for the resources
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
# "staging-klb-backup-2017.05.28.1930-test"
#
# The resource group name pattern is important to be know since
# you must manage these resource groups and delete them.
# There will be also patterns on how snapshots are stored inside
# this resource group but they are not documented and you should not
# rely on them, they are implementation details.
#
# Backup resource groups should never be changed, because of that
# they are read only locked after all snapshots are added. There is
# also a delete lock to avoid deleting backups on accident.
# The azure_vm_backup_delete function will release the locks and delete
# a backup resource group for you.
#
# During the backup procedure the VM will be shutdown, and restarted
# after all snapshots are taken.
fn azure_vm_backup_create(vmname, group, prefix) {

}
