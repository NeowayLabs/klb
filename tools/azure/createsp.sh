#!/usr/bin/env nash

ROLES = (
	("Owner" "true")
	("Reader" "true")
)

fn check_role(role) {
	for entry in $ROLES {
		if $entry[0] == $role {
			return $entry[1]
		}
	}

	return "false"
}

defRole = "Owner"

arglen    <= len($ARGS)
_, status <= test $arglen -lt 4

if $status == "0" {
	print("Usage: %s <subscription_name> <service_principal> <password> [<role>]\n", $ARGS[0])
	exit("1")
}

subscription_name = $ARGS[1]
service_principal = $ARGS[2]
password          = $ARGS[3]
role              = $defRole

if $arglen == "5" {
	role = $ARGS[4]
}
if check_role($role) == "false" {
	print("Invalid role:[%s] - possible roles are: Owner or Reader\n", $role)
	exit("1")
}

print("Setting Azure subscription: subscription_name[%s]\n", $subscription_name)

out, status <= az account set --subscription $subscription_name >[2=1]

if $status != "0" {
	print("Unable to setting Azure subscription, error[%s]\n", $out)
	exit("1")
}

print("Getting Azure subscription id: subscription_name[%s]\n", $subscription_name)

out, status <= (
	az account show
		--subscription $subscription_name
		--output=json | jq -r ".id" >[2=1]
)

if $status != "0" {
	print("Unable to getting Azure subscription id, error[%s]\n", $out)
	exit("1")
} else {
	subscription_id = $out
}

print("Creating Azure Active Directory: service principal[%s]\n", $service_principal)

out, status <= (
	az ad sp create-for-rbac
	    --name $service_principal
	    --password $password
	    --role $role
	    --scope "/subscriptions/"+$subscription_id >[2=1]
)

if $status != "0" {
	print("Unable to create Azure Active Directory service principal, error[%s]\n", $out)
	exit("1")
}
