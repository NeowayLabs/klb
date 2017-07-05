# Droplet related functions

fn digitalocean_droplet_new(name, size, image, region) {
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

fn digitalocean_droplet_set_backups(instance) {
	instance <= append($instance, "--enable-backups")

	return $instance
}

fn digitalocean_droplet_set_ipv6(instance) {
	instance <= append($instance, "--enable-ipv6")

	return $instance
}

fn digitalocean_droplet_set_monitoring(instance) {
	instance <= append($instance, "--enable-monitoring")

	return $instance
}

fn digitalocean_droplet_set_private_network(instance) {
	instance <= append($instance, "--enable-private-networking")

	return $instance
}

fn digitalocean_droplet_set_ssh_key(instance, keys) {
	instance <= append($instance, "--ssh-keys")
	instance <= append($instance, $keys)

	return $instance
}

fn digitalocean_droplet_set_user_data_file(instance, file) {
	instance <= append($instance, "--user-data-file")
	instance <= append($instance, $file)

	return $instance
}

fn digitalocean_droplet_set_tag_name(instance, tag) {
	instance <= append($instance, "--tag-name")
	instance <= append($instance, $tag)

	return $instance
}

fn digitalocean_droplet_set_tag_names(instance, tags) {
	instance <= append($instance, "--tag-names")
	instance <= append($instance, $tags)

	return $instance
}

fn digitalocean_droplet_set_tag_volumes(instance, volumes) {
	instance <= append($instance, "--volumes")
	instance <= append($instance, $volumes)

	return $instance
}

fn digitalocean_droplet_set_wait(instance) {
	instance <= append($instance, "--wait")

	return $instance
}

fn digitalocean_droplet_create(instance) {
	doctl compute droplet create $instance
}
