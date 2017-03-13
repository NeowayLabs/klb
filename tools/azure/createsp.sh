#!/usr/bin/env nash

if len($ARGS) != "4" {
	echo "usage: "+$ARGS[0]+" <subscription_id>  <service_principal> <password>"
	exit
}

subscription_id   = $ARGS[1]
service_principal = $ARGS[2]
password          = $ARGS[3]

subscription_name <= azure account show $subscription_id --json | jq -r ".[0].name"

echo "Setting azure subscription to: "+$subscription_name+" ["+$subscription_id+"]"
azure account set $subscription_id >[1=]
echo "Creating azure service principal: "+$service_principal

service_principal_obj_id <= azure ad sp create -n $service_principal -p $password --json | jq -r ".objectId"

fn loop(func, count) {
        sequence <= seq 1 $count
        range <= split($sequence, "\n")
        for i in $range {
            $func()
        }
}

fn grant_access() {
        echo "Granting access to "+$service_principal+" ["+$service_principal_obj_id+"] at: "+$subscription_name+" ["+$subscription_id+"]"
        -azure role assignment create --objectId $service_principal_obj_id -o Owner -c "/subscriptions/"+$subscription_id+"/" >[1=]
        if $status == "0" {
            echo "granted access with success"
            exit("0")
        }
        echo "unable to grant access, trying again"
        sleep 1
}

loop($grant_access, "5")

echo "unable to grant access, try again expired"
exit("1")
