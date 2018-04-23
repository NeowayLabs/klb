# Logging in with a service principal
#
# Service principals are like user accounts to which you can apply rules using
# Azure Active Directory. Authenticating with a service principal is the best way
# to secure the usage of your Azure resources from either your scripts or
# applications that manipulate resources.
#
# Login, should be called before executing any functions
#
# Login info will be loaded from:
# AZURE_TENANT_ID
# AZURE_CLIENT_ID
# AZURE_CLIENT_SECRET
# AZURE_SERVICE_PRINCIPAL
# AZURE_SUBSCRIPTION_ID

fn azure_login() {
	subscriptionID = $AZURE_SUBSCRIPTION_ID
	tenantID       = $AZURE_TENANT_ID
	clientID       = $AZURE_CLIENT_ID
	secretID       = $AZURE_CLIENT_SECRET
	username       = $AZURE_SERVICE_PRINCIPAL
	
	azure_login_credentials($subscriptionID, $tenantID, $clientID, $secretID, $username)
}

fn azure_login_credentials(subscriptionID, tenantID, clientID, secretID, username) {
	# azure cli 2.0
	az account clear
	az login --service-principal -u $username -p $secretID --tenant $tenantID --output table
	az account set --subscription $subscriptionID

	# azure cli 1.0
	azure telemetry --disable
	azure config mode arm
	azure login -q -u $clientID --service-principal --tenant $tenantID -p $secretID
	azure account set $subscriptionID
}
