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

# Azure CLI 1.0
fn azure_login() {
        azure config mode arm

        tenantID = $AZURE_TENANT_ID
        clientID = $AZURE_CLIENT_ID
        secretID = $AZURE_CLIENT_SECRET

        azure login -q -u $clientID --service-principal --tenant $tenantID -p $secretID
}

# Azure CLI 2.0
fn az_login() {

        username = $AZURE_SERVICE_PRINCIPAL
        tenantID = $AZURE_TENANT_ID
        secretID = $AZURE_CLIENT_SECRET

        az account clear
        az login --service-principal -u $username -p $AZURE_CLIENT_SECRET --tenant $tenantID --output table
}
