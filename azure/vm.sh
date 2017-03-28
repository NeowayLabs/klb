# Machine related functions

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

fn azure_vm_set_vmsize(instance, vmsize) {
	instance <= append($instance, "--vm-size")
	instance <= append($instance, $vmsize)

	return $instance
}

fn azure_vm_set_customdata(instance, customdata) {
	instance <= append($instance, "--custom-data")
	instance <= append($instance, $customdata)

	return $instance
}

fn azure_vm_set_username(instance, username) {
	instance <= append($instance, "--admin-username")
	instance <= append($instance, $username)

	return $instance
}

fn azure_vm_set_publickeyfile(instance, publickeyfile) {
	instance <= append($instance, "--ssh-publickey-file")
	instance <= append($instance, $publickeyfile)

	return $instance
}

fn azure_vm_set_availset(instance, availset) {
	instance <= append($instance, "--availset-name")
	instance <= append($instance, $availset)

	return $instance
}

fn azure_vm_set_vnet(instance, vnet) {
	instance <= append($instance, "--vnet-name")
	instance <= append($instance, $vnet)

	return $instance
}

fn azure_vm_set_subnet(instance, subnet) {
	instance <= append($instance, "--vnet-subnet-name")
	instance <= append($instance, $subnet)

	return $instance
}

fn azure_vm_set_nic(instance, nic) {
	instance <= append($instance, "--nic-name")
	instance <= append($instance, $nic)

	return $instance
}

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
	instance <= append($instance, "--nic-names")
	instance <= append($instance, $nics)

	return $instance
}

fn azure_vm_set_storageaccount(instance, storageaccount) {
	instance <= append($instance, "--storage-account-name")
	instance <= append($instance, $storageaccount)

	return $instance
}

fn azure_vm_set_osdiskvhd(instance, osdiskvhd) {
	instance <= append($instance, "--os-disk-vhd")
	instance <= append($instance, $osdiskvhd)

	return $instance
}

fn azure_vm_set_datadiskvhd(instance, datadiskvhd) {
	instance <= append($instance, "--data-disk-vhd")
	instance <= append($instance, $datadiskvhd)

	return $instance
}

fn azure_vm_set_datadisksize(instance, datadisksize) {
	instance <= append($instance, "--data-disk-size")
	instance <= append($instance, $datadisksize)

	return $instance
}

fn azure_vm_set_imageurn(instance, imageurn) {
	instance <= append($instance, "--image-urn")
	instance <= append($instance, $imageurn)

	return $instance
}

# azure_vm_set_disablebootdiagnostic switch parameter to disable
# the boot diagnostics (by default enabled) for VM.
# `instance` is the instance of vm.
fn azure_vm_set_disablebootdiagnostics(instance) {
	instance <= append($instance, "--disable-boot-diagnostics")

	return $instance
}

# azure_vm_set_bootdiagnosticsstorageuri specify storage account for boot diagnostics.
# `instance` is the instance of vm.
# `storageaccount` is the name of storage account.
fn azure_vm_set_bootdiagnosticsstorage(instance, storageaccount) {
	instance <= append($instance, "--boot-diagnostics-storage-uri")
	instance <= append($instance, "https://"+$storageaccount+".blob.core.windows.net/")

	return $instance
}

fn azure_vm_create(instance) {
	azure vm create $instance
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
