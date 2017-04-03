# Droplet related functions

fn digitalocean_droplet_create(name, size, image, region) {
	instance = (
		$name
		"--size"
		$size
		"--image"
		$image
		"--region"
		$region
	)

	return $instance
}

fn digitalocean_vm_set_backups(instance) {
	instance <= append($instance, "--enable-backups")

	return $instance
}

fn digitalocean_vm_set_ipv6(instance) {
	instance <= append($instance, "--enable-ipv6")

	return $instance
}

fn digitalocean_vm_set_monitoring(instance) {
	instance <= append($instance, "--enable-monitoring")

	return $instance
}

fn digitalocean_vm_set_private_network(instance) {
	instance <= append($instance, "--enable-private-networking")

	return $instance
}

fn digitalocean_vm_set_ssh_key(instance, keys) {
	instance <= append($instance, "--ssh-keys")
	instance <= append($instance, $keys)

	return $instance
}

fn digitalocean_droplet_create(instance) {
	doctl compute droplet create $instance
}

# TODO
#      --tag-name string             Tag name
#      --tag-names value             Tag names (default [])
#      --user-data-file string       User data file
#      --volumes value               Volumes to attach (default [])
#      --wait                        Wait for droplet to be created
