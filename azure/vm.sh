# Machine related functions

# azure_vm_new creates a new instance of "virtual machine".
# `name` is the name of the virtual machine.
# `group` is name of resource group.
# `location` is the Azure Region.
# `ostype` is the type of OS installed on a custom VHD.
# Accepted values: linux, windows
fn azure_vm_new(name, group, location, ostype) {
	instance = (
		"--name"
		$name
		"--resource-group"
		$group
		"--location"
		$location
		"--os-type"
		$ostype
	)

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
fn azure_vm_set_nic(instance, nicnames) {
	instance <= append($instance, "--nics")
	instance <= append($instance, $nics)

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
# `storagesku` is the the sku of storage account to persist VM. By default, only
# Standard_LRS and Premium_LRS are allowed.
fn azure_vm_set_storagesku(instance, storagesku) {
	instance <= append($instance, "--storage-sku")
	instance <= append($instance, $imageurn)

	return $instance
}

# azure_vm_create creates an Azure "Virtual Machine".
# `instance` is the name of the instance.
fn azure_vm_create(instance) {
	az vm create $instance
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
