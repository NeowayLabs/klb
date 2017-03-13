#!/usr/bin/env nash

if len($ARGS) != 4 {
        echo "usage: " + $ARGS[0] + " <subscription_id> <service_principal> <password>"
        exit("0")
}

subscription_id = $ARGS[1]
service_principal = $ARGS[2]
password = $ARGS[3]

azure account set $subscription_id
azure ad sp create -n $service_principal -p $password
