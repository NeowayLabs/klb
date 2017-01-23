package azure

import (
	"errors"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type VM struct {
	client compute.VirtualMachinesClient
	f      fixture.F
}

func NewVM(f fixture.F) *VM {
	as := &VM{
		client: compute.NewVirtualMachinesClient(f.Session.SubscriptionID),
		f:      f,
	}
	as.client.Authorizer = f.Session.Token
	return as
}

// AssertExists checks if VM exists in the resource group.
// Fail tests otherwise.
func (vm *VM) AssertExists(t *testing.T, name, expectedAvailSet, expectedVMSize, expectedOsType, expectedNic string) {
	vm.f.Retrier.Run(newID("VM", "AssertExists", name), func() error {
		v, err := vm.client.Get(vm.f.ResGroupName, name, "")
		if err != nil {
			return err
		}
		if v.VirtualMachineProperties == nil {
			return errors.New("Field VirtualMachineProperties is nil!")
		}
		properties := *v.VirtualMachineProperties
		if properties.AvailabilitySet == nil {
			return errors.New("Field AvailabilitySet is nil!")
		}
		if properties.AvailabilitySet.ID == nil {
			return errors.New("Field ID is nil!")
		}
		gotAvailSet := *properties.AvailabilitySet.ID
		if !strings.Contains(gotAvailSet, strings.ToUpper(expectedAvailSet)) {
			return errors.New("AvailSet expected is " + expectedAvailSet + " but got " + gotAvailSet)
		}
		if properties.HardwareProfile == nil {
			return errors.New("Field HardwareProfile is nil!")
		}
		hardwareProfile := *properties.HardwareProfile
		gotVMSize := string(hardwareProfile.VMSize)
		if gotVMSize != expectedVMSize {
			return errors.New("VM Size expected is " + expectedVMSize + " but got " + gotVMSize)
		}
		if properties.StorageProfile == nil {
			return errors.New("Field StorageProfile is nil!")
		}
		if properties.StorageProfile.OsDisk == nil {
			return errors.New("Field OsDisk is nil!")
		}
		osDisk := *properties.StorageProfile.OsDisk
		gotOsType := string(osDisk.OsType)
		if gotOsType != expectedOsType {
			return errors.New("OS type expected is " + expectedOsType + " but got " + gotOsType)
		}
		if properties.NetworkProfile == nil {
			return errors.New("Field NetworkProfile is nil!")
		}
		network := *properties.NetworkProfile.NetworkInterfaces
		if len(network) == 0 {
			return errors.New("Field NetworkInterfaces is nil!")
		}
		net := network[0]
		if net.ID == nil {
			return errors.New("Field ID is nil!")
		}
		gotNic := string(*net.ID)
		if !strings.Contains(gotNic, expectedNic) {
			return errors.New("Nic expected is " + expectedNic + " but got " + gotNic)
		}
		return nil
	})
}
