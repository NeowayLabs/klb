package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
)

type ResourceGroup struct {
	client resources.Client
}

func NewResourceGroup(t *testing.T, s *Session) *ResourceGroup {
	rg := &ResourceGroup{client: resources.NewClient(s.subscriptionID)}
	rg.client.Authorizer = s.token
	return rg
}
