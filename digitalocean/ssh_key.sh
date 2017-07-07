# ssh key related functions

# digitalocean_ssh_key_create create an
# ssh key into the Digital Ocean platform and
# returns all the ssh key information in JSON format.
# If an error occurs, it returns an empty string and
# the error.
#
# `name` is the ssh key name.
# `key` is the ssh key contents to create.
fn digitalocean_ssh_key_create(name, key) {
	out, status <= doctl compute ssh-key create $name --public-key $key --output json
	if $status != "0" {
	   return "", $out
	}

	return $out, ""
}

# digitalocean_ssh_key_import imports an existent
# ssh key into the Digital Ocean platform and
# returns all the ssh key information in JSON format.
# If an error occurs, it returns an empty string and
# the error.
#
# `name` is the ssh key name.
# `key` is the ssh key file to import.
fn digitalocean_ssh_key_import(name, key) {
	out, status <= doctl compute ssh-key import $name --public-key-file $key --output json
	if $status != "0" {
	   return "", $out
	}

	return $out, ""
}

fn digitalocean_ssh_key_exists(name) {
	key, status <= (
		doctl compute ssh-key list --output json |
		jq ".[].id" |
		grep $name
	)

	if $status == "0" {
		return "0"
	}

	return "1"
}
