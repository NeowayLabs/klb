package azure

import (
	"errors"
	"fmt"
	"strings"
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
func (nic *Nic) AssertExists(t *testing.T, name string, privateIP string) {
	nic.f.Retrier.Run(newID("Nic", "AssertExists", name), func() error {

		ipconfigs, err := nic.GetIPConfigs(t, name)
		if err != nil {
			return err
		}

		for _, ipconfig := range ipconfigs {
			if ipconfig.PrivateIPAddress == privateIP {
				return nil
			}
		}

		return fmt.Errorf("unable to find privateIP[%s] on ipconfigs[%s]", privateIP, ipconfigs)
	})
}

type NicIPConfig struct {
	Name                  string
	PrivateIPAddress      string
	LBBackendAddrPoolsIDs []string
}

func parseNICID(t *testing.T, nicID string) (string, string) {
	// example:
	// /subscriptions/e3e74e5f-cc81-49d1-8fab-00fff864c080/resourceGroups/klb-VMGetIPAddress-1508439367-7943/providers/Microsoft.Network/networkInterfaces/nic1
	// FIXME: Did not found a way to get directly by the ID =(
	// This is a bad/coupled solution
	parsed := strings.Split(nicID, "/")
	if len(parsed) != 9 {
		t.Fatalf("unexpected NIC ID[%s] parsed[%s], M$ protocol probably changed", nicID, parsed)
	}
	return parsed[4], parsed[8]
}

func (nic *Nic) GetIPConfigsByID(t *testing.T, ID string) ([]NicIPConfig, error) {
	resgroup, nicname := parseNICID(t, ID)
	return nic.getIPConfigs(t, resgroup, nicname)
}

func (nic *Nic) GetIPConfigs(t *testing.T, name string) ([]NicIPConfig, error) {
	return nic.getIPConfigs(t, nic.f.ResGroupName, name)
}

func (nic *Nic) getIPConfigs(t *testing.T, resgroup string, nicname string) ([]NicIPConfig, error) {

	var ipconfigs []NicIPConfig

	wraperror := func(err error) error {
		return fmt.Errorf("Nic.GetInfo: error[%s]", err)
	}

	n, err := nic.client.Get(resgroup, nicname, "")
	if err != nil {
		return []NicIPConfig{}, wraperror(err)
	}

	if n.InterfacePropertiesFormat == nil {
		return []NicIPConfig{}, wraperror(errors.New("The field InterfacePropertiesFormat is nil!"))
	}

	propertiesFormat := *n.InterfacePropertiesFormat

	if propertiesFormat.IPConfigurations == nil {
		return []NicIPConfig{}, wraperror(errors.New("No IPConfigurations found on NIC"))
	}

	for _, azIPConfig := range *propertiesFormat.IPConfigurations {
		ipconfig := NicIPConfig{}

		if azIPConfig.Name == nil {
			return []NicIPConfig{}, wraperror(fmt.Errorf("ip config[%s] has nil Name", azIPConfig))
		}

		if azIPConfig.InterfaceIPConfigurationPropertiesFormat == nil {
			return []NicIPConfig{}, wraperror(fmt.Errorf("ip config[%s] has nil properties", azIPConfig))
		}

		if azIPConfig.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress == nil {
			return []NicIPConfig{}, wraperror(fmt.Errorf("ip config[%s] has nil private IP address", azIPConfig))
		}

		ipconfig.Name = *azIPConfig.Name
		ipconfig.PrivateIPAddress = *azIPConfig.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress

		if azIPConfig.InterfaceIPConfigurationPropertiesFormat.LoadBalancerBackendAddressPools != nil {
			pools := *azIPConfig.InterfaceIPConfigurationPropertiesFormat.LoadBalancerBackendAddressPools
			poolsIDs := []string{}

			for _, pool := range pools {
				if pool.ID == nil {
					return []NicIPConfig{}, wraperror(fmt.Errorf("pool[%s] has no name", pool))
				}
				poolsIDs = append(poolsIDs, *pool.ID)
			}

			ipconfig.LBBackendAddrPoolsIDs = poolsIDs
		}

		ipconfigs = append(ipconfigs, ipconfig)
	}

	return ipconfigs, nil
}
