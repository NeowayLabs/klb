package azure

import (
	"os"
	"testing"

	restazure "github.com/Azure/go-autorest/autorest/azure"
)

type Session struct {
	ClientID       string
	ClientSecret   string
	SubscriptionID string
	TenantID       string
	token          *restazure.ServicePrincipalToken
}

func (s *Session) generateToken(t *testing.T) {
	oauthConfig, err := restazure.PublicCloud.OAuthConfigForTenant(s.TenantID)
	if err != nil {
		t.Fatal(err)
	}
	s.token, err = restazure.NewServicePrincipalToken(
		*oauthConfig,
		s.ClientID,
		s.ClientSecret,
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
		ClientID:       getenv(t, "AZURE_CLIENT_ID"),
		ClientSecret:   getenv(t, "AZURE_CLIENT_SECRET"),
		SubscriptionID: getenv(t, "AZURE_SUBSCRIPTION_ID"),
		TenantID:       getenv(t, "AZURE_TENANT_ID"),
	}
	session.generateToken(t)
	return session
}
