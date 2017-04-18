package fixture

import (
	"fmt"
	"os"
	"testing"

	restazure "github.com/Azure/go-autorest/autorest/azure"
)

type Session struct {
	ClientID         string
	ClientSecret     string
	SubscriptionID   string
	TenantID         string
	ServicePrincipal string
	Token            *restazure.ServicePrincipalToken
}

// Env provides all environment variables required to
// run scripts that will integrate with Azure.
func (s *Session) Env() []string {
	return []string{
		fmt.Sprintf("AZURE_CLIENT_ID=%s", s.ClientID),
		fmt.Sprintf("AZURE_CLIENT_SECRET=%s", s.ClientSecret),
		fmt.Sprintf("AZURE_SUBSCRIPTION_ID=%s", s.SubscriptionID),
		fmt.Sprintf("AZURE_TENANT_ID=%s", s.TenantID),
		fmt.Sprintf("AZURE_SERVICE_PRINCIPAL=%s", s.ServicePrincipal),
	}
}

func (s *Session) generateToken(t *testing.T) {
	oauthConfig, err := restazure.PublicCloud.OAuthConfigForTenant(s.TenantID)
	if err != nil {
		t.Fatal(err)
	}
	s.Token, err = restazure.NewServicePrincipalToken(
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
		ClientID:         getenv(t, "AZURE_CLIENT_ID"),
		ClientSecret:     getenv(t, "AZURE_CLIENT_SECRET"),
		SubscriptionID:   getenv(t, "AZURE_SUBSCRIPTION_ID"),
		TenantID:         getenv(t, "AZURE_TENANT_ID"),
		ServicePrincipal: getenv(t, "AZURE_SERVICE_PRINCIPAL"),
	}
	session.generateToken(t)
	return session
}
