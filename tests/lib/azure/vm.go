package azure

import (
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
func (vm *VM) AssertExists(t *testing.T, name string) {
	vm.f.Retrier.Run(newID("VM", "AssertExists", name), func() error {
		_, err := vm.client.Get(vm.f.ResGroupName, name)
		return err
	})
}
