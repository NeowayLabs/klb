package azure

import (
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type PublicIp struct {
	client network.PublicIPAddressesClient
	f      fixture.F
}

func NewPublicIp(f fixture.F) *PublicIp {
	as := &PublicIp{
		client: network.NewPublicIPAddressesClient(f.Session.SubscriptionID),
		f:      f,
	}
	as.client.Authorizer = f.Session.Token
	return as
}

// AssertExists checks if publicIp exists in the resource group.
// Fail tests otherwise.
func (publicIp *PublicIp) AssertExists(t *testing.T, name string) {
	publicIp.f.Retrier.Run(newID("PublicIp", "AssertExists", name), func() error {
		n, err := publicIp.client.Get(publicIp.f.ResGroupName, name, "")
		if err != nil {
			return err
		}

		if n.PublicIPAddressPropertiesFormat == nil {
			return errors.New("The field PublicIPAddressPropertiesFormat is nil!")
		}
		properties := *n.PublicIPAddressPropertiesFormat

		if properties.IPAddress == nil {
			return errors.New("The field IPAddress is nil!")
		}

		return nil
	})
}
