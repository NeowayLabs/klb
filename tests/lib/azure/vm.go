package azure

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type VM struct {
	client compute.VirtualMachinesClient
	f      fixture.F
}

type VMDataDisk struct {
	Lun                int
	Name               string
	SizeGB             int
	Caching            string
	StorageAccountType string
}

type VMOsDisk struct {
	OsType  string
	Name    string
	SizeGB  int
	Caching string
}

func NewVM(f fixture.F) *VM {
	as := &VM{
		client: compute.NewVirtualMachinesClient(f.Session.SubscriptionID),
		f:      f,
	}
	as.client.Authorizer = f.Session.Token
	return as
}

func (vm *VM) IPs(t *testing.T, vmname string) []string {
	abortonerr := func(err error) {
		if err != nil {
			t.Fatalf("error[%s] getting IP address for vm[%s]", err, vmname)
		}
	}

	v, err := vm.client.Get(vm.f.ResGroupName, vmname, "")
	abortonerr(err)

	if v.VirtualMachineProperties == nil {
		abortonerr(errors.New("no virtual machine properties found"))
	}
	properties := v.VirtualMachineProperties

	if properties.NetworkProfile == nil {
		abortonerr(errors.New("Field NetworkProfile is nil!"))
	}
	network := *properties.NetworkProfile.NetworkInterfaces
	if len(network) == 0 {
		abortonerr(errors.New("Field NetworkInterfaces is nil!"))
	}

	ips := []string{}
	nic := NewNic(vm.f)

	for _, net := range network {
		if net.ID == nil {
			abortonerr(errors.New("network interface ID is nil!"))
		}

		ipConfigs, err := nic.GetIPConfigsByID(t, *net.ID)
		abortonerr(err)

		for _, ipConfig := range ipConfigs {
			ips = append(ips, ipConfig.PrivateIPAddress)
		}
	}

	return ips
}

func (vm *VM) OsDisk(t *testing.T, vmname string) VMOsDisk {

	var osdisk *VMOsDisk

	vm.f.Retrier.Run(newID("VM", "OsDisk", vmname), func() error {
		v, err := vm.client.Get(vm.f.ResGroupName, vmname, "")
		if err != nil {
			return err
		}
		if v.VirtualMachineProperties == nil {
			return fmt.Errorf("no virtual machine properties found on vm %s", vmname)
		}
		if v.VirtualMachineProperties.StorageProfile == nil {
			return fmt.Errorf("no storage profile found on vm %s", vmname)
		}

		storageProfile := v.VirtualMachineProperties.StorageProfile
		if storageProfile.OsDisk == nil {
			return fmt.Errorf("no os disk found on vm %s", vmname)
		}

		if storageProfile.OsDisk.Name == nil {
			return errors.New("os disk has no name")
		}

		if storageProfile.OsDisk.DiskSizeGB == nil {
			return errors.New("os disk has no size")
		}

		osdisk = &VMOsDisk{
			Name:    *storageProfile.OsDisk.Name,
			SizeGB:  int(*storageProfile.OsDisk.DiskSizeGB),
			OsType:  string(storageProfile.OsDisk.OsType),
			Caching: string(storageProfile.OsDisk.Caching),
		}

		return nil
	})

	if osdisk == nil {
		t.Fatal("unable to get os disks for vm %q", vmname)
	}

	return *osdisk
}

func (vm *VM) DataDisks(t *testing.T, vmname string) []VMDataDisk {

	disksinfo := []VMDataDisk{}

	vm.f.Retrier.Run(newID("VM", "DataDisks", vmname), func() error {
		v, err := vm.client.Get(vm.f.ResGroupName, vmname, "")
		if err != nil {
			return err
		}
		if v.VirtualMachineProperties == nil {
			return fmt.Errorf("no virtual machine properties found on vm %s", vmname)
		}
		if v.VirtualMachineProperties.StorageProfile == nil {
			return fmt.Errorf("no storage profile found on vm %s", vmname)
		}

		storageProfile := v.VirtualMachineProperties.StorageProfile
		if storageProfile.DataDisks == nil {
			return fmt.Errorf("no data disks found on vm %s", vmname)
		}

		for _, disk := range *storageProfile.DataDisks {
			if disk.Name == nil {
				continue
			}
			if disk.Lun == nil {
				continue
			}
			if disk.DiskSizeGB == nil {
				continue
			}
			if disk.ManagedDisk == nil {
				continue
			}
			disksinfo = append(disksinfo, VMDataDisk{
				Name:               *disk.Name,
				Lun:                int(*disk.Lun),
				SizeGB:             int(*disk.DiskSizeGB),
				Caching:            string(disk.Caching),
				StorageAccountType: string(disk.ManagedDisk.StorageAccountType),
			})
		}

		return nil
	})

	return disksinfo
}

// AssertAttachedDisk checks if VM has the following disk attached
func (vm *VM) AssertAttachedDataDisk(
	t *testing.T,
	vmname string,
	diskname string,
	diskSizeGB int,
	storageAccountType string,
	caching string,
) {
	vm.f.Retrier.Run(newID("VM", "AssertAttachedDataDisk", vmname), func() error {

		disks := vm.DataDisks(t, vmname)

		for _, disk := range disks {
			vm.f.Logger.Printf("got disk %+v", disk)

			if disk.Name != diskname {
				continue
			}
			if disk.SizeGB != diskSizeGB {
				continue
			}
			if disk.StorageAccountType != storageAccountType {
				continue
			}
			if disk.Caching != caching {
				continue
			}
			return nil
		}

		return fmt.Errorf("unable to find disk %q on vm %q. available disks: %s", diskname, vmname, disks)
	})
}

// AssertExistsByName checks if VM exists in the resource group
// based only on its name. Fail tests otherwise.
func (vm *VM) AssertExistsByName(t *testing.T, name string) {
	vm.f.Retrier.Run(newID("VM", "AssertExistsByName", name), func() error {
		_, err := vm.client.Get(vm.f.ResGroupName, name, "")
		if err != nil {
			return fmt.Errorf("unable to find vm %q, error: %s", name, err)
		}
		return nil
	})
}

// AssertExists checks if VM exists in the resource group.
// Fail tests otherwise.
func (vm *VM) AssertExists(
	t *testing.T,
	name string,
	expectedAvailSet string,
	expectedVMSize string,
	expectedNic string,
	expectedTags string,
) {
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
		tagSplit := strings.Split(expectedTags, "=")
		value, ok := (*v.Tags)[tagSplit[0]]
		if !ok {
			return errors.New("Tag " + tagSplit[0] + " not found")
		}
		gotTag := string(*value)
		if gotTag != tagSplit[1] {
			return errors.New("Tag value expected is " + tagSplit[1] + " but got " + gotTag)
		}
		return nil
	})
}
