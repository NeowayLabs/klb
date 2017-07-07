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

# digitalocean_ssh_key_exists checks if a
# ssh key exists in the Digital Ocean platform and
# returns the ssh key ID.
# If an error occurs, it returns an empty string.
#
# `name` is the ssh key name.
fn digitalocean_ssh_key_exists(name) {
	out, status <= (
		doctl compute ssh-key list --output json |
		jq "map(select(.name==\""+$name+"\"))| .[].id"
	)
	if $status != "0" {
	   return ""
	}

	return $out
}
