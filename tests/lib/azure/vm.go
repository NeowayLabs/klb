package azure

import (
	"errors"
	"fmt"
	"log"
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
func (vm *VM) AssertExists(t *testing.T, name, expectedAvailSet, expectedVMSize, expectedOsType string) {
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
		hardwareProfile := *properties.HardwareProfile
		gotVMSize := string(hardwareProfile.VMSize)
		if gotVMSize != expectedVMSize {
			return errors.New("VM Size expected is " + expectedVMSize + " but got " + gotVMSize)
		}
		osDisk := *properties.StorageProfile.OsDisk
		gotOsType := string(osDisk.OsType)
		if gotOsType != expectedOsType {
			return errors.New("OS type expected is " + expectedOsType + " but got " + gotOsType)
		}
		fmt.Println("STORAGE PROFILE: %v", *properties.StorageProfile, *properties.StorageProfile.ImageReference, *properties.StorageProfile.OsDisk, *properties.StorageProfile.DataDisks, "OS PROFILE %v", *properties.OsProfile.ComputerName)
		network := *properties.NetworkProfile.NetworkInterfaces
		log.Fatal("NETWORK PROFILE:%v", *properties.NetworkProfile, *network[0].ID)
		return nil
	})
}
