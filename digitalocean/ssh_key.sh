# ssh key related functions

fn digitalocean_ssh_key_create(name, key) {
	doctl compute ssh-key create $name --public-key $key
}

fn digitalocean_ssh_key_import(name, key) {
	doctl compute ssh-key import $name --public-key-file $key
}

fn digitalocean_ssh_key_exists(name) {
	key, status <= (
		doctl compute ssh-key list --output json |
		jq ".[].name" |
		grep $name
	)

	if $status == "0" {
		return "0"
	}

	return "1"
}
