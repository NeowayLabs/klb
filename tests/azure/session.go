package azure

import (
	"os"
	"testing"

	restazure "github.com/Azure/go-autorest/autorest/azure"
)

type Session struct {
	clientID       string
	clientSecret   string
	subscriptionID string
	tenantID       string
	token          *restazure.ServicePrincipalToken
}

func (s *Session) generateToken(t *testing.T) {
	oauthConfig, err := restazure.PublicCloud.OAuthConfigForTenant(s.tenantID)
	if err != nil {
		t.Fatal(err)
	}
	s.token, err = restazure.NewServicePrincipalToken(
		*oauthConfig,
		s.clientID,
		s.clientSecret,
		restazure.PublicCloud.ResourceManagerEndpoint,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func getenv(t *testing.T, varname string) string {
	val := os.Getenv(varname)
	if val == "" {
		t.Fatalf("Missing environment variable %q", varname)
	}
	return val
}

func NewSession(t *testing.T) *Session {
	session := &Session{
		clientID:       getenv(t, "AZURE_CLIENT_ID"),
		clientSecret:   getenv(t, "AZURE_CLIENT_SECRET"),
		subscriptionID: getenv(t, "AZURE_SUBSCRIPTION_ID"),
		tenantID:       getenv(t, "AZURE_TENANT_ID"),
	}
	session.generateToken(t)
	return session
}
