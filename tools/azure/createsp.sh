#!/usr/bin/env nash

if len($ARGS) != "4" {
	echo "usage: "+$ARGS[0]+" <service_principal>  <subscription_id> <password>"
	exit
}

SERVICE_PRINCIPAL = $ARGS[1]
SUBSCRIPTION_ID   = $ARGS[2]
PASSWORD          = $ARGS[3]

SUBSCRIPTION_NAME <= azure account show $SUBSCRIPTION_ID --json | jq -r ".[0].name"

echo "Setting azure subscription to: "+$SUBSCRIPTION_NAME+" ["+$SUBSCRIPTION_ID+"]"
azure account set $SUBSCRIPTION_ID >[1=]
echo "Creating azure service principal: "+$SERVICE_PRINCIPAL

SERVICE_PRINCIPAL_OBJECT_ID <= azure ad sp create -n $SERVICE_PRINCIPAL -p $PASSWORD --json | jq -r ".objectId"

echo "Granting access to "+$SERVICE_PRINCIPAL+" ["+$SERVICE_PRINCIPAL_OBJECT_ID+"] at: "+$SUBSCRIPTION_NAME+" ["+$SUBSCRIPTION_ID+"]"
sleep 5
azure role assignment create --objectId $SERVICE_PRINCIPAL_OBJECT_ID -o Owner -c "/subscriptions/"+$SUBSCRIPTION_ID+"/" >[1=]
echo "Done"
