package azure_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
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

func createVM(t *testing.T, f fixture.F, vmSize string, sku string) string {
	vm := genVMName()
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
		sku,
	)

	f.Logger.Println("creating VM")
	vms := azure.NewVM(f)
	vms.AssertExists(t, vm, resources.availSet, vmSize, resources.nic)
	f.Logger.Println("created VM with success, attaching a disk")

	return vm
}

func testVMCreation(t *testing.T, f fixture.F, vmSize string, sku string) {

	vm := createVM(t, f, vmSize, sku)

	diskname := "createVMExtraDisk"
	size := 10

	attachNewDiskOnVM(t, f, vm, diskname, size, sku)

	vms := azure.NewVM(f)
	vms.AssertAttachedDataDisk(t, vm, diskname, size, sku)

	f.Logger.Println("VM with attached disk created with success")
}

func testStandardDiskVM(t *testing.T, f fixture.F) {
	testVMCreation(t, f, "Basic_A0", "Standard_LRS")
}

func testPremiumDiskVM(t *testing.T, f fixture.F) {
	testVMCreation(t, f, "Standard_DS4_v2", "Premium_LRS")
}

func testVMSnapshot(t *testing.T, f fixture.F, vmSize string, sku string) {
	vm := createVM(t, f, vmSize, sku)

	disks := []struct {
		name string
		size int
	}{
		{name: "disk1", size: 10},
		{name: "disk2", size: 20},
		{name: "disk3", size: 30},
	}

	vms := azure.NewVM(f)

	for _, disk := range disks {
		attachNewDiskOnVM(t, f, vm, disk.name, disk.size, sku)
		vms.AssertAttachedDataDisk(t, vm, disk.name, disk.size, sku)
	}

	outfile, err := ioutil.TempFile("", "create_vm_snapshots_output")
	if err != nil {
		t.Fatalf("error creating output file: %s", err)
	}
	defer os.Remove(outfile.Name()) // clean up

	f.Shell.Run("./testdata/create_vm_snapshots.sh", f.ResGroupName, vm, outfile.Name())

	f.Logger.Println("created snapshots, retrieving ids")
	idsraw, err := ioutil.ReadAll(outfile)
	if err != nil {
		t.Fatalf("error reading output file: %s", err)
	}

	ids := strings.Split(string(idsraw), "\n")
	f.Logger.Printf("parsed ids: %s", ids)

	if len(ids) != len(disks) {
		t.Fatalf("expected %d snapshots, got %d", len(disks), len(ids))
	}

	vmbackup := createVM(t, f, vmSize, sku)
	for _, id := range ids {
		// attachDiskOnVM(t, f, vmbackup, id)
		// TODO: Check the attached disk properties
	}
}

func testVMSnapshotStandard(t *testing.T, f fixture.F) {
	testVMSnapshot(t, f, "Basic_A0", "Standard_LRS")
}

func testVMSnapshotPremium(t *testing.T, f fixture.F) {
	testVMSnapshot(t, f, "Standard_DS4_v2", "Premium_LRS")
}

func attachNewDiskOnVM(
	t *testing.T,
	f fixture.F,
	vmname string,
	diskname string,
	diskSizeGB int,
	sku string,
) {
	f.Shell.Run(
		"./testdata/attach_new_disk.sh",
		f.ResGroupName,
		vmname,
		diskname,
		strconv.Itoa(diskSizeGB),
		sku,
	)
}

func attachDiskOnVM(t *testing.T, f fixture.F, vmname string, diskid string) {
	f.Shell.Run("./testdata/attach_disk.sh", f.ResGroupName, vmname, diskid)
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
	fixture.Run(t, "VMCreationStandardDisk", 30*time.Minute, location, testStandardDiskVM)
	fixture.Run(t, "VMCreationPremiumDisk", 30*time.Minute, location, testPremiumDiskVM)
	fixture.Run(t, "VMSnapshotStandard", 30*time.Minute, location, testVMSnapshotStandard)
	fixture.Run(t, "VMSnapshotPremium", 30*time.Minute, location, testVMSnapshotPremium)
}
