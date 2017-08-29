package azure

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type Nic struct {
	client network.InterfacesClient
	f      fixture.F
}

func NewNic(f fixture.F) *Nic {
	as := &Nic{
		client: network.NewInterfacesClient(f.Session.SubscriptionID),
		f:      f,
	}
	as.client.Authorizer = f.Session.Token
	return as
}

// AssertExists checks if nic exists in the resource group.
// Fail tests otherwise.
func (nic *Nic) AssertExists(t *testing.T, name string, address string) {
	nic.f.Retrier.Run(newID("Nic", "AssertExists", name), func() error {
		// FIXME: Use GetIPConfigs here to check existence of NIC

		n, err := nic.client.Get(nic.f.ResGroupName, name, "")
		if err != nil {
			return err
		}

		if n.InterfacePropertiesFormat == nil {
			return errors.New("The field InterfacePropertiesFormat is nil!")
		}
		properties := *n.InterfacePropertiesFormat

		if properties.IPConfigurations == nil {
			return errors.New("The field IPConfigurations is nil!")
		}
		ip := *properties.IPConfigurations

		if len(ip) == 0 || ip[0].PrivateIPAddress == nil {
			return errors.New("The field PrivateIPAddress is nil!")
		}

		privateAddress := *ip[0].PrivateIPAddress
		if privateAddress != address {
			return errors.New("Nic created with wrong Address. Expected: " + address + "Actual: " + privateAddress)
		}

		return nil
	})
}

type NicIPConfig struct {
	Name                            string
	PrivateIPAddress                string
	LoadBalancerBackendAddressPools []string
}

func (nic *Nic) GetIPConfigs(t *testing.T, name string, address string) ([]NicIPConfig, error) {

	var properties []NicIPConfig

	wraperror := func(err error) error {
		return fmt.Errorf("Nic.GetInfo: error[%s]", err)
	}

	n, err := nic.client.Get(nic.f.ResGroupName, name, "")
	if err != nil {
		return properties, wraperror(err)
	}

	if n.InterfacePropertiesFormat == nil {
		return properties, wraperror(errors.New("The field InterfacePropertiesFormat is nil!"))
	}

	propertiesFormat := *n.InterfacePropertiesFormat

	if propertiesFormat.IPConfigurations == nil {
		return properties, wraperror(errors.New("No IPConfigurations found on NIC"))
	}

	//ipconfigs := *propertiesFormat.IPConfigurations
	//for _, ipconfig := range ipconfigs {
	//}

	return properties, nil
}
