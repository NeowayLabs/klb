# Scaleset related functions

# azure_vmss_new creates a new instance of "Virtual Machines Scale Set".
# `name` is the name of the scale set.
# `group` is name of resource group.
# `location` is the Azure Region.
fn azure_vmss_new(name, group, location) {
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

# azure_vmss_set_imageurn sets the image of "Virtual Machines Scale Set".
# `instance` is the name of the instance.
# `imageurn` is the name of the operating system image (URN alias, URN,
# Custom Image name or ID, or VHD Blob URI)
fn azure_vmss_set_imageurn(instance, imageurn) {
	instance <= append($instance, "--image")
	instance <= append($instance, $imageurn)

	return $instance
}

# azure_vmss_set_publickeyfile sets the SSH public key of "Virtual Machines Scale Set".
# `instance` is the name of the instance.
# `publickeyfile` is the SSH public key or public key file path.
fn azure_vmss_set_publickeyfile(instance, publickeyfile) {
	instance <= append($instance, "--ssh-key-value")
	instance <= append($instance, $publickeyfile)

	return $instance
}

# azure_vmss_set_customdata sets the custom init script file of "Virtual Machines Scale Set".
# `instance` is the name of the instance.
# `customdata` is the location of custom init script file or text (cloud-config).
fn azure_vmss_set_customdata(instance, customdata) {
	instance <= append($instance, "--custom-data")
	instance <= append($instance, $customdata)

	return $instance
}

# azure_vmss_set_vmsize sets the size or sku of "Virtual Machine" in "Scale Set".
# `instance` is the name of the instance.
# `size` or `sku` is the size of virtual machine.
# See https://goo.gl/vxW7we for size info.
fn azure_vmss_set_vmsize(instance, vmsize) {
	instance <= append($instance, "--vm-sku")
	instance <= append($instance, $vmsize)

	return $instance
}

# azure_vmss_set_ostype sets ostype of "Virtual Machines Scale Set".
# `instance` is the name of the instance.
# `ostype` is the type of OS installed on a custom VHD.
fn azure_vmss_set_ostype(instance, ostype) {
	instance <= append($instance, "--os-type")
	instance <= append($instance, $ostype)

	return $instance
}

# azure_vmss_set_instancecount sets instance count of "Virtual Machines Scale Set".
# `instance` is the name of the instance.
# `count` is a number of VM's instances on "scale set".
fn azure_vmss_set_instancecount(instance, count) {
	instance <= append($instance, "--instance-count")
	instance <= append($instance, $count)

	return $instance
}

fn azure_vmss_set_vnet(instance, vnet) {
	instance <= append($instance, "--vnet-name")
	instance <= append($instance, $vnet)

	return $instance
}

fn azure_vmss_set_subnet(instance, subnet) {
	instance <= append($instance, "--subnet")
	instance <= append($instance, $subnet)

	return $instance
}

# azure_vmss_set_lb sets instance a name to use when creating a new load balancer (default) or referencing an existing one.
# `instance` is the name of the instance.
# `lb` reference an existing load balancer by ID or specify "" for none.
fn azure_vmss_set_lb(instance, lb) {
	instance <= append($instance, "--lb")
	instance <= append($instance, $lb)

	return $instance
}

# azure_vm_create creates a "Virtual Machine Scale Set".
# `instance` is the name of the instance.
fn azure_vmss_create(instance) {
	az vmss create $instance
}