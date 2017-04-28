# Disk related functions

# azure_disk_new creates a new instance of "managed disk".
# `name` is the name of the managed disk.
# `group` is name of resource group.
# `location` is the Azure Region.
fn azure_disk_new(name, group, location) {
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

# azure_disk_set_size sets the size of "managed disk".
# `instance` is the name of the instance.
# `size` is the size in Gb of managed disk.
fn azure_disk_set_size(instance, size) {
	instance <= append($instance, "--size-gb")
	instance <= append($instance, $size)

	return $instance
}

# azure_disk_set_sku sets the kind of "managed disk".
# `instance` is the disk instance.
# `sku` is the underlying storage sku. Allowed values: Premium_LRS,
#  Standard_LRS. Default: Premium_LRS.
fn azure_disk_set_sku(instance, sku) {
	instance <= append($instance, "--sku")
	instance <= append($instance, $sku)

	return $instance
}

# azure_disk_set_source sets the source of "managed disk".
# `instance` is the disk instance.
# `source` is the source to create the disk from, including a sas uri
# to a blob, managed disk id or name, or snapshot id or name.
fn azure_disk_set_source(instance, source) {
	instance <= append($instance, "--source")
	instance <= append($instance, $source)

	return $instance
}

# azure_disk_new creates a new "managed disk".
# `instance` is the disk instance.
fn azure_disk_create(instance) {
	az disk create --output table $instance
}

# azure_disk_get_id returns the id of a previously created
# disk. This id is used to attach the disk on a VM.
fn azure_disk_get_id(resgroup, name) {
	res <= az disk show -g $resgroup -n $name --query "id"

	return $res
}
