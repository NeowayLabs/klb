package azure

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/disk"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type Disks struct {
	f      fixture.F
	client disk.DisksClient
}

func NewDisk(f fixture.F) *Disks {
	as := &Disks{
		f:      f,
		client: disk.NewDisksClient(f.Session.SubscriptionID),
	}
	as.client.Authorizer = f.Session.Token
	return as
}

// AssertExists checks if disk exists in the resource group.
// Fail tests otherwise.
func (d *Disks) AssertExists(t *testing.T, name string, size string, sku string) {
	d.f.Retrier.Run(newID("Disk", "AssertExists", name), func() error {
		res, err := d.client.Get(d.f.ResGroupName, name)
		if err != nil {
			return err
		}
		if res.Properties == nil {
			return errors.New("no properties found on disk")
		}
		gotSKU := string(res.Properties.AccountType)
		if gotSKU != sku {
			return fmt.Errorf(
				"expected type %q got %q",
				sku,
				gotSKU,
			)
		}
		d.f.Logger.Println(res.Properties.AccountType)

		wantSizeInt, err := strconv.Atoi(size)
		if err != nil {
			return err
		}
		wantSize := int32(wantSizeInt)
		gotSize := *res.Properties.DiskSizeGB
		if wantSize != gotSize {
			return fmt.Errorf(
				"want disksize %d but got %d",
				wantSize,
				gotSize,
			)
		}
		return nil
	})
}
