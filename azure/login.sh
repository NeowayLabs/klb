# Login, should be called before executing any functions
# Login info will be loaded from:
# AZURE_TENANT_ID
# AZURE_CLIENT_ID
# AZURE_CLIENT_SECRET

fn azure_login() {
        azure config mode arm

        tenantID = $AZURE_TENANT_ID
        clientID = $AZURE_CLIENT_ID
        secretID = $AZURE_CLIENT_SECRET

        azure login -q -u $clientID --service-principal --tenant $tenantID -p $secretID
}
