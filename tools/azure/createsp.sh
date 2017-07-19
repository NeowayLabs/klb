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

arglen  <= len($ARGS)
_, cond <= test $arglen -lt 4

if $cond == "0" {
	echo "usage: "+$ARGS[0]+" <subscription_id> <service_principal> <password> [<role>]"
	exit
}

subscription_id   = $ARGS[1]
service_principal = $ARGS[2]
password          = $ARGS[3]
role              = $defRole

if $arglen == "5" {
	role = $ARGS[4]
}
if check_role($role) == "false" {
	echo "Invalid role:"+$role+" - possible roles are: Owner or Reader"
	
	exit("1")
}

subscription_name <= azure account show $subscription_id --json | jq -r ".[0].name"

echo "Setting azure subscription to: "+$subscription_name+" ["+$subscription_id+"]"
azure account set $subscription_id >[1=]
echo "Creating azure service principal: "+$service_principal

service_principal_obj_id <= azure ad sp create -n $service_principal -p $password --json | jq -r ".objectId"

fn loop(func, count) {
	sequence <= seq 1 $count
	range    <= split($sequence, "\n")

	for i in $range {
		$func()
	}
}

fn grant_access() {
	echo "Granting access to "+$service_principal+" ["+$service_principal_obj_id+"] at: "+$subscription_name+" ["+$subscription_id+"]"
	-azure role assignment create --objectId $service_principal_obj_id -o $role -c "/subscriptions/"+$subscription_id+"/" >[1=]

	if $status == "0" {
		echo "granted access with success with role: "+$role
		
		exit("0")
	}

	echo "unable to grant access, trying again with role: "+$role
	sleep 1
}

loop($grant_access, "5")

echo "unable to grant access, try again expired"

exit("1")
