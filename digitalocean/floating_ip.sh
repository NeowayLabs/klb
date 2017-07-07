# Floating IP related functions

# digitalocean_floating_ip_create creates a "Floating IP" and
# returns the IPv4 information of the "Floating IP".
# If an error occurs, it returns an empty string and
# the error.
#
# `region` is the floating ip region to be created.
fn digitalocean_floating_ip_create(region) {
	out, status <= (
		doctl compute floating-ip create --region $region --output json |
		jq ".[].ip" |
		tr -d "\""
	)

	if $status != "0" {
		return "", $out
	}

	return $out, ""
}

# digitalocean_floating_ip_droplet_create creates
# a "Floating IP" assigned to a "Droplet" and returns
# the IPv4 information of the "Floating IP".
# If an error occurs, it returns an empty string and
# the error.
#
# `droplet_id` is the floating ip droplet id to be created.
fn digitalocean_floating_ip_droplet_create(droplet_id) {
	out, status <= (
		doctl compute floating-ip create --droplet-id $droplet_id --output json |
		jq ".[].ip" |
		tr -d "\""
	)

	if $status != "0" {
		return "", $out
	}

	return $out, ""
}

# digitalocean_floating_ip_droplet_assign assigns
# a "Floating IP" to a "Droplet".
#
# `floating_ip` is the floating ip to be assigned.
# `droplet_id` is the droplet id to be assigned.
fn digitalocean_floating_ip_droplet_assign(floating_ip, droplet_id) {
	(
		doctl compute floating-ip-action assign $floating_ip $droplet_id
	)
}
