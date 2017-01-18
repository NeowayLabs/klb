#!/usr/bin/env nash


if len($ARGS) != "4" {
    echo "Usage: " $ARGS[0] "<(sh|nash)> <service principal name>" "<service secret>"
	abort
}

shell=$ARGS[1]

fn printvar(name, value) {
        if $shell == "nash" {
	    printf "setenv %s=\"%s\"\n" $name $value
	    return
        }
	printf "export %s=\"%s\"\n" $name $value
}

SPNAME                  = $ARGS[2]
SPSECRET                = $ARGS[3]

AZURE_SUBSCRIPTION_ID   <= (
	azure account show
			--json
			 |
	jq -r ".[0].id" |
	tr -d "\n"
)

AZURE_SUBSCRIPTION_NAME <= (
	azure account show
			--json
			 |
	jq -r ".[0].name" |
	tr -d "\n"
)

AZURE_TENANT_ID         <= (
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
