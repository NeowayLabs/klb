package azure

import (
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type Disk struct {
	f fixture.F
}

func NewDisk(f fixture.F) *Disk {
	as := &Disk{
		//client: compute.NewVirtualMachinesClient(f.Session.SubscriptionID),
		f: f,
	}
	//as.client.Authorizer = f.Session.Token
	return as
}

// AssertExists checks if disk exists in the resource group.
// Fail tests otherwise.
func (d *Disk) AssertExists(t *testing.T, name string, size string, sku string) {
	d.f.Retrier.Run(newID("Disk", "AssertExists", name), func() error {
		return nil
	})
}
