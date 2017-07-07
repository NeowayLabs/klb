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

fn digitalocean_droplet_set_ssh_key(instance, key) {
	instance <= append($instance, "--ssh-keys")
	instance <= append($instance, $key)

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

# digitalocean_droplet_create creates a "Droplet" and
# returns all the "Droplet" information in JSON format.
# If an error occurs, it returns an empty string and
# the error.
#
# `instance` is the droplet parameters instance.
fn digitalocean_droplet_create(instance) {
	out, status <= doctl compute droplet create $instance --output json
	if $status != "0" {
	   return "", $out
	}
	
	return $out, ""
}

# digitalocean_droplet_exists checks if a
# "Droplet" exists in the Digital Ocean platform and
# returns the "Droplet" ID.
# If an error occurs, it returns an empty string.
#
# `name` is the droplet name.
fn digitalocean_droplet_exists(name) {
	out, status <= (
		doctl compute droplet list --output json |
		jq "map(select(.name==\""+$name+"\"))| .[].id"
	)
	if $status != "0" {
	   return ""
	}

	return $out
}

# digitalocean_droplet_get_ip gets the public IPv4
# associated with the "Droplet".
# If an error occurs, it returns an empty string.
#
# `name` is the droplet name.
fn digitalocean_droplet_get_ip(name) {
	out, status <= (
		doctl compute droplet list --output json |
		jq "map(select(.name==\""+$name+"\"))| .[].networks.v4[].ip_address" |
		tr -d "\""
	)
	if $status != "0" {
	   return ""
	}

	return $out
}
