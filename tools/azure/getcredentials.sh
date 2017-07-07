#!/usr/bin/env nash

if len($ARGS) != "5" {
	echo "Usage: " $ARGS[0] "<(sh|nash)> <subscription name> <service principal name> <service secret>"
	
	exit("0")
}

shell = $ARGS[1]

fn printvar(name, value) {
	if $shell == "nash" {
		printf "setenv %s=\"%s\"\n" $name $value
		
		return
	}

	printf "export %s=\"%s\"\n" $name $value
}

AZURE_SUBSCRIPTION_NAME = $ARGS[2]
SPNAME                  = $ARGS[3]
SPSECRET                = $ARGS[4]

azure login
azure account set $AZURE_SUBSCRIPTION_NAME > /dev/null

AZURE_SUBSCRIPTION_ID <= (
	azure account show
			--json
			 |
	jq -r ".[0].id" |
	tr -d "\n"
)

AZURE_TENANT_ID       <= (
	azure account show
			--json
			 |
	jq -r ".[0].tenantId" |
	tr -d "\n"
)

printvar("AZURE_SUBSCRIPTION_ID", $AZURE_SUBSCRIPTION_ID)
printvar("AZURE_SUBSCRIPTION_NAME", $AZURE_SUBSCRIPTION_NAME)
printvar("AZURE_TENANT_ID", $AZURE_TENANT_ID)

AZURE_CLIENT_ID <= (
	azure ad sp show -c $SPNAME --json |
	jq -r ".[0].appId" |
	tr -d "\n"
)

printvar("AZURE_CLIENT_ID", $AZURE_CLIENT_ID)
printvar("AZURE_CLIENT_SECRET", $SPSECRET)
printvar("AZURE_SERVICE_PRINCIPAL", "http://"+$SPNAME)
