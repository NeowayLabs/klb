package azure_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genVMName() string {
	return fmt.Sprintf("klbvmtests%d", rand.Intn(1000))
}

type VMResources struct {
	availSet string
	vnet     string
	subnet   string
	nic      string
}

func testVMCreate(t *testing.T, f fixture.F) {

	vm := genVMName()
	vmSize := "Standard_DS4_v2"
	username := "core"
	osDisk := "test.vhd"
	imageUrn := "OpenLogic:CentOS:7.2:7.2.20161026"
	keyFile := "./testdata/key.pub"

	resources := createVMResources(t, f)

	f.Shell.Run(
		"./testdata/create_vm.sh",
		vm,
		f.ResGroupName,
		f.Location,
		vmSize,
		username,
		resources.availSet,
		resources.nic,
		osDisk,
		imageUrn,
		keyFile,
	)

	f.Logger.Println("creating VM")
	vms := azure.NewVM(f)
	vms.AssertExists(t, vm, resources.availSet, vmSize, resources.nic)

	f.Logger.Println("created VM with success, attaching a disk")
	diskname := "createVMExtraDisk"
	size := 10
	sku := "Standard_LRS"
	createDisk(t, f, diskname, size, sku)
	f.Logger.Println("created disk, attaching it")

	attachDiskOnVM(t, f, vm, diskname)
	f.Logger.Println("VM with attached disk created with success")
}

func attachDiskOnVM(t *testing.T, f fixture.F, vmname string, diskname string) {
	// TODO: Improve attached disk validation
	//f.Shell.Run(
	//"./testdata/attach_disk.sh",
	//f.ResGroupName,
	//vmname,
	//diskname,
	//)
}

func createVMResources(t *testing.T, f fixture.F) VMResources {

	resources := VMResources{}
	resources.availSet = genAvailSetName()
	resources.vnet = genVnetName()
	resources.subnet = genSubnetName()
	resources.nic = genNicName()

	nsg := genNsgName()
	vnetAddress := "10.116.0.0/16"
	subnetAddress := "10.116.1.0/24"
	addrnic := "10.116.1.100"
	updatedomain := "3"
	faultdomain := "3"

	f.Shell.Run(
		"./testdata/create_vm_avail_set.sh",
		f.ResGroupName,
		resources.availSet,
		f.Location,
		updatedomain,
		faultdomain,
	)

	createVNet(t, f, vnetDescription{
		name:     resources.vnet,
		vnetAddr: vnetAddress,
	})

	f.Shell.Run(
		"./testdata/create_nsg.sh",
		nsg,
		f.ResGroupName,
		f.Location,
	)

	f.Shell.Run(
		"./testdata/create_subnet.sh",
		resources.subnet,
		f.ResGroupName,
		resources.vnet,
		subnetAddress,
		nsg,
	)

	f.Shell.Run(
		"./testdata/create_nic.sh",
		f.ResGroupName,
		resources.nic,
		f.Location,
		resources.vnet,
		resources.subnet,
		addrnic,
	)

	return resources
}

func TestVM(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "CreateVM", 25*time.Minute, location, testVMCreate)
}
